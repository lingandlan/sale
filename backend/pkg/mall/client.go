package mall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"marketplace/backend/internal/config"
)

// WSYClient WSY商城HTTP客户端
type WSYClient struct {
	baseURL    string
	appID      string
	appSecret  string
	customerID string
	httpClient *http.Client

	mu          sync.Mutex
	accessToken string
	tokenExpire time.Time
}

// NewWSYClient 创建WSY客户端
func NewWSYClient(cfg config.MallConfig) *WSYClient {
	return &WSYClient{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		appID:      cfg.AppID,
		appSecret:  cfg.AppSecret,
		customerID: cfg.CustomerID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// wsyCommonResp WSY通用响应
type wsyCommonResp struct {
	ErrCode int             `json:"errcode"`
	ErrMsg  string           `json:"errmsg"`
	Data    json.RawMessage  `json:"data"`
}

// GetAccessToken 获取access_token（自动缓存，提前5分钟刷新）
func (c *WSYClient) GetAccessToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.tokenExpire.Add(-5*time.Minute)) {
		return c.accessToken, nil
	}

	endpoint := fmt.Sprintf("%s/wsy_pub/third_api/index.php?m=third_application_authorization&a=get_access_token", c.baseURL)
	form := url.Values{
		"appid":             {c.appID},
		"appsecret":         {c.appSecret},
		"customer_id_lock":  {c.customerID},
	}

	resp, err := c.httpClient.PostForm(endpoint, form)
	if err != nil {
		return "", fmt.Errorf("WSY get_access_token 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("WSY get_access_token 读取响应失败: %w", err)
	}

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Data    struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("WSY get_access_token 解析失败: %w, body: %s", err, string(body))
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("WSY get_access_token 失败: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	c.accessToken = result.Data.AccessToken
	expiresIn := time.Duration(result.Data.ExpiresIn) * time.Second
	if expiresIn == 0 {
		expiresIn = 2 * time.Hour
	}
	c.tokenExpire = time.Now().Add(expiresIn)

	return c.accessToken, nil
}

// PhoneToUserID 通过手机号获取商城用户ID
func (c *WSYClient) PhoneToUserID(phone string) (string, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/wsy_pub/third_api/index.php?m=third_application_authorization&a=third_function", c.baseURL)
	form := url.Values{
		"access_token":     {token},
		"act":              {"10000_phone_get_user_info"},
		"customer_id_lock": {c.customerID},
		"country_code":     {"+86"},
		"phone":            {phone},
	}

	body, err := c.postForm(endpoint, form)
	if err != nil {
		return "", err
	}

	var result struct {
		ErrCode  int    `json:"errcode"`
		ErrMsg   string `json:"errmsg"`
		UserInfo struct {
			UserID string `json:"user_id"`
		} `json:"user_info"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("WSY phone_get_user_info 解析失败: %w, body: %s", err, string(body))
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("WSY phone_get_user_info 失败: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}
	if result.UserInfo.UserID == "" {
		return "", fmt.Errorf("该手机号未注册商城会员")
	}

	return result.UserInfo.UserID, nil
}

// GetUserIntegral 查询用户积分余额
func (c *WSYClient) GetUserIntegral(userID string) (float64, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return 0, err
	}

	endpoint := fmt.Sprintf("%s/wsy_pub/third_api/index.php?m=third_application_authorization&a=third_function", c.baseURL)
	form := url.Values{
		"access_token":     {token},
		"act":              {"10000_integral_user_integral"},
		"customer_id_lock": {c.customerID},
		"user_id":          {userID},
	}

	body, err := c.postForm(endpoint, form)
	if err != nil {
		return 0, err
	}

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Data    struct {
			UserID   string `json:"user_id"`
			Integral string `json:"integral"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("WSY integral_user_integral 解析失败: %w, body: %s", err, string(body))
	}
	if result.ErrCode != 0 {
		return 0, fmt.Errorf("WSY integral_user_integral 失败: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	var integral float64
	fmt.Sscanf(result.Data.Integral, "%f", &integral)
	return integral, nil
}

// postForm 通用POST表单请求
func (c *WSYClient) postForm(endpoint string, form url.Values) ([]byte, error) {
	resp, err := c.httpClient.PostForm(endpoint, form)
	if err != nil {
		return nil, fmt.Errorf("WSY 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("WSY 读取响应失败: %w", err)
	}
	return body, nil
}
