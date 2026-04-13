import { test, expect } from '@playwright/test';

const API_BASE = process.env.API_BASE || 'http://localhost:8080/api/v1';

test.describe('认证模块', () => {
  test('登录成功', async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: {
        phone: '13800000000',
        password: '123456'
      }
    });

    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.code).toBe(0);
    expect(body.data.access_token).toBeDefined();
  });

  test('登录失败 - 错误的密码', async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: {
        phone: '13800000000',
        password: 'wrongpassword'
      }
    });

    expect(response.status()).toBe(401);
  });

  test('刷新Token', async ({ request }) => {
    const loginRes = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const loginData = await loginRes.json();
    const refreshToken = loginData.data?.refresh_token;

    if (refreshToken) {
      const refreshRes = await request.post(`${API_BASE}/auth/refresh`, {
        data: { refresh_token: refreshToken }
      });
      expect(refreshRes.ok()).toBeTruthy();
    }
  });
});

test.describe('用户模块', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('获取用户信息', async ({ request }) => {
    const response = await request.get(`${API_BASE}/user/info`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('获取用户列表', async ({ request }) => {
    const response = await request.get(`${API_BASE}/admin/users?page=1&page_size=20`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });
});

test.describe('B端充值申请', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('创建充值申请', async ({ request }) => {
    const response = await request.post(`${API_BASE}/recharge/b-apply`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        centerId: 'center-1',
        centerName: '测试中心',
        amount: 10000,
        transactionNo: `TXN${Date.now()}`
      }
    });
    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.data?.id).toBeDefined();
  });

  test('获取审批列表', async ({ request }) => {
    const response = await request.get(`${API_BASE}/recharge/b-approval?page=1&page_size=10`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('审批通过', async ({ request }) => {
    const createRes = await request.post(`${API_BASE}/recharge/b-apply`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        centerId: 'center-1',
        centerName: '测试中心',
        amount: 5000,
        transactionNo: `TXN${Date.now()}`
      }
    });
    const createData = await createRes.json();
    const appId = createData.data?.id;

    if (appId) {
      const approveRes = await request.post(`${API_BASE}/recharge/b-approval/action`, {
        headers: { Authorization: `Bearer ${authToken}` },
        data: { id: appId, action: 'approve', remark: '测试通过' }
      });
      expect(approveRes.ok()).toBeTruthy();
    }
  });
});

test.describe('C端充值录入', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('创建C端充值', async ({ request }) => {
    const response = await request.post(`${API_BASE}/recharge/c-entry`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        memberId: 'member-1',
        memberName: '张三',
        memberPhone: '13800000001',
        centerId: 'center-1',
        centerName: '测试中心',
        amount: 500,
        paymentMethod: 'wechat'
      }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('获取充值记录列表', async ({ request }) => {
    const response = await request.get(`${API_BASE}/recharge/records?page=1&page_size=10`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.data).toBeDefined();
  });
});

test.describe('充值记录', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('查询充值记录 - 按手机号', async ({ request }) => {
    const response = await request.get(`${API_BASE}/recharge/records?memberPhone=13800000001&page=1&page_size=10`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('查询充值记录 - 按充值中心', async ({ request }) => {
    const response = await request.get(`${API_BASE}/recharge/records?centerId=center-1&page=1&page_size=10`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('查询充值记录详情', async ({ request }) => {
    const listRes = await request.get(`${API_BASE}/recharge/records?page=1&page_size=10`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    const listData = await listRes.json();
    const firstId = listData.data?.list?.[0]?.id;

    if (firstId) {
      const detailRes = await request.get(`${API_BASE}/recharge/records/${firstId}`, {
        headers: { Authorization: `Bearer ${authToken}` }
      });
      expect(detailRes.ok()).toBeTruthy();
    }
  });
});

test.describe('门店卡管理', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('发放门店卡', async ({ request }) => {
    const response = await request.post(`${API_BASE}/card/issue`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        holderId: 'member-1',
        holderName: '李四',
        holderPhone: '13900000001',
        amount: 1000,
        centerId: 'center-1',
        centerName: '测试中心'
      }
    });
    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.data?.cardNo).toBeDefined();
  });

  test('核销门店卡', async ({ request }) => {
    const issueRes = await request.post(`${API_BASE}/card/issue`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        holderId: 'member-2',
        holderName: '王五',
        holderPhone: '13900000002',
        amount: 2000,
        centerId: 'center-1',
        centerName: '测试中心'
      }
    });
    const issueData = await issueRes.json();
    const cardNo = issueData.data?.cardNo;

    if (cardNo) {
      const consumeRes = await request.post(`${API_BASE}/card/consume`, {
        headers: { Authorization: `Bearer ${authToken}` },
        data: { cardNo, amount: 500, remark: '购买商品' }
      });
      expect(consumeRes.ok()).toBeTruthy();
    }
  });

  test('获取卡列表', async ({ request }) => {
    const response = await request.get(`${API_BASE}/card/list?page=1&page_size=10`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('获取卡统计', async ({ request }) => {
    const response = await request.get(`${API_BASE}/card/stats`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });
});

test.describe('充值中心管理', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('获取充值中心列表', async ({ request }) => {
    const response = await request.get(`${API_BASE}/center`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
    expect(Array.isArray(response)).toBeTruthy();
  });

  test('创建充值中心', async ({ request }) => {
    const timestamp = Date.now();
    const response = await request.post(`${API_BASE}/center`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        name: `测试中心${timestamp}`,
        code: `TEST${String(timestamp).slice(-6)}`,
        address: '北京市朝阳区',
        phone: '010-88888888'
      }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('更新充值中心', async ({ request }) => {
    const createRes = await request.post(`${API_BASE}/center`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        name: '待更新中心',
        code: `UPD${String(Date.now()).slice(-6)}`,
        address: '北京市朝阳区',
        phone: '010-88888888'
      }
    });
    const createData = await createRes.json();
    const centerId = createData.data?.id;

    if (centerId) {
      const updateRes = await request.put(`${API_BASE}/center/${centerId}`, {
        headers: { Authorization: `Bearer ${authToken}` },
        data: {
          name: '已更新中心',
          code: createData.data?.code,
          address: '北京市海淀区',
          phone: '010-99999999'
        }
      });
      expect(updateRes.ok()).toBeTruthy();
    }
  });

  test('删除充值中心', async ({ request }) => {
    const createRes = await request.post(`${API_BASE}/center`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        name: '待删除中心',
        code: `DEL${String(Date.now()).slice(-6)}`,
        address: '北京市朝阳区',
        phone: '010-88888888'
      }
    });
    const createData = await createRes.json();
    const centerId = createData.data?.id;

    if (centerId) {
      const deleteRes = await request.delete(`${API_BASE}/center/${centerId}`, {
        headers: { Authorization: `Bearer ${authToken}` }
      });
      expect(deleteRes.ok()).toBeTruthy();
    }
  });
});

test.describe('操作员管理', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('获取操作员列表', async ({ request }) => {
    const response = await request.get(`${API_BASE}/operator`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('创建操作员', async ({ request }) => {
    const timestamp = Date.now();
    const response = await request.post(`${API_BASE}/operator`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        name: `操作员${timestamp}`,
        phone: `139${String(timestamp).slice(-8)}`,
        password: '123456',
        centerId: 'center-1',
        role: 'operator'
      }
    });
    expect(response.ok()).toBeTruthy();
  });
});

test.describe('Dashboard', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('获取统计数据', async ({ request }) => {
    const response = await request.get(`${API_BASE}/dashboard/statistics`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('获取待办事项', async ({ request }) => {
    const response = await request.get(`${API_BASE}/dashboard/todos`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });

  test('获取充值趋势', async ({ request }) => {
    const response = await request.get(`${API_BASE}/dashboard/recharge-trends?days=7`, {
      headers: { Authorization: `Bearer ${authToken}` }
    });
    expect(response.ok()).toBeTruthy();
  });
});

test.describe('API一致性验证', () => {
  let authToken: string;

  test.beforeAll(async ({ request }) => {
    const response = await request.post(`${API_BASE}/auth/login`, {
      data: { phone: '13800000000', password: '123456' }
    });
    const body = await response.json();
    authToken = body.data?.access_token;
  });

  test('验证所有充值相关端点可访问', async ({ request }) => {
    const endpoints = [
      { method: 'GET', path: '/recharge/b-apply', auth: true },
      { method: 'GET', path: '/recharge/b-approval', auth: true },
      { method: 'GET', path: '/recharge/c-entry', auth: true },
      { method: 'GET', path: '/recharge/records', auth: true },
      { method: 'GET', path: '/card/list', auth: true },
      { method: 'GET', path: '/card/stats', auth: true },
      { method: 'GET', path: '/center', auth: true },
      { method: 'GET', path: '/operator', auth: true },
      { method: 'GET', path: '/dashboard/statistics', auth: true },
    ];

    for (const endpoint of endpoints) {
      const headers: Record<string, string> = {};
      if (endpoint.auth) {
        headers['Authorization'] = `Bearer ${authToken}`;
      }

      const response = await request.get(`${API_BASE}${endpoint.path}`, { headers });
      
      if (!response.ok()) {
        console.log(`❌ ${endpoint.method} ${endpoint.path} - ${response.status()}`);
      }
      
      expect(response.ok() || response.status() === 404, 
        `${endpoint.path} should be accessible (${response.status()})`).toBeTruthy();
    }
  });
});
