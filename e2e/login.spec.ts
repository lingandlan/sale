import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:5175';

test.describe('登录页面 E2E 测试', () => {
  
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.waitForLoadState('networkidle');
  });

  test('TC-001: 空表单提交', async ({ page }) => {
    await page.click('button:has-text("登录")');
    // 应该有验证提示
    await expect(page.locator('text=/请输入|必填|不能为空/i')).toBeVisible({ timeout: 3000 });
  });

  test('TC-002: 只填手机号提交', async ({ page }) => {
    await page.fill('input[placeholder*="手机号"]', '13800138000');
    await page.click('button:has-text("登录")');
    await expect(page.locator('text=/请输入.*密码|密码.*必填/i')).toBeVisible({ timeout: 3000 });
  });

  test('TC-003: 只填密码提交', async ({ page }) => {
    await page.fill('input[placeholder*="密码"]', 'password123');
    await page.click('button:has-text("登录")');
    await expect(page.locator('text=/请输入.*手机|手机.*必填/i')).toBeVisible({ timeout: 3000 });
  });

  test('TC-004: 错误密码', async ({ page }) => {
    await page.fill('input[placeholder*="手机号"]', '13800138000');
    await page.fill('input[placeholder*="密码"]', 'wrongpassword');
    await page.click('button:has-text("登录")');
    // 应该显示错误 toast
    await expect(page.locator('.u-toast-text, [class*="toast"], text=/密码错误|账号不存在/i')).toBeVisible({ timeout: 5000 });
  });

  test('TC-005: 正确登录 (需真实账号)', async ({ page }) => {
    // TODO: 填入真实测试账号
    const TEST_PHONE = process.env.TEST_PHONE || 'YOUR_TEST_PHONE';
    const TEST_PASSWORD = process.env.TEST_PASSWORD || 'YOUR_TEST_PASSWORD';
    
    if (TEST_PHONE === 'YOUR_TEST_PHONE') {
      test.skip('需要配置测试账号');
    }

    await page.fill('input[placeholder*="手机号"]', TEST_PHONE);
    await page.fill('input[placeholder*="密码"]', TEST_PASSWORD);
    await page.click('button:has-text("登录")');
    
    // 验证跳转到首页
    await expect(page).toHaveURL(/\/$|\/home|\/dashboard/, { timeout: 10000 });
  });

  test('TC-006: 记住密码复选框', async ({ page }) => {
    const checkbox = page.locator('input[type="checkbox"]');
    await expect(checkbox).not.toBeChecked();
    await checkbox.check();
    await expect(checkbox).toBeChecked();
    await checkbox.uncheck();
    await expect(checkbox).not.toBeChecked();
  });

  test('TC-007: 忘记密码链接', async ({ page }) => {
    await page.click('text=忘记密码？');
    // 验证跳转到忘记密码页面
    await expect(page).toHaveURL(/forgot|reset|password/, { timeout: 5000 });
  });
});
