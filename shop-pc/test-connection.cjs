#!/usr/bin/env node

/**
 * 前后端联调测试脚本
 * 测试PC端前端与后端API的集成
 */

const http = require('http');

// 颜色输出
const colors = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  blue: '\x1b[36m'
};

function log(message, color = 'reset') {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

function makeRequest(options, data = null) {
  return new Promise((resolve, reject) => {
    const req = http.request(options, (res) => {
      let body = '';
      res.on('data', chunk => body += chunk);
      res.on('end', () => {
        try {
          resolve({
            statusCode: res.statusCode,
            headers: res.headers,
            body: JSON.parse(body)
          });
        } catch (e) {
          resolve({
            statusCode: res.statusCode,
            headers: res.headers,
            body: body
          });
        }
      });
    });

    req.on('error', reject);

    if (data) {
      req.write(JSON.stringify(data));
    }

    req.end();
  });
}

async function testBackendHealth() {
  log('\n========================================', 'blue');
  log('测试1: 后端服务健康检查', 'blue');
  log('========================================', 'blue');

  try {
    const response = await makeRequest({
      hostname: 'localhost',
      port: 8080,
      path: '/health',
      method: 'GET'
    });

    if (response.statusCode === 200 && response.body.status === 'ok') {
      log('✓ 后端服务正常运行', 'green');
      log(`  响应: ${JSON.stringify(response.body)}`, 'green');
      return true;
    } else {
      log(`✗ 后端服务异常: ${response.statusCode}`, 'red');
      return false;
    }
  } catch (error) {
    log(`✗ 无法连接到后端服务: ${error.message}`, 'red');
    return false;
  }
}

async function testLoginAPI() {
  log('\n========================================', 'blue');
  log('测试2: 登录API接口', 'blue');
  log('========================================', 'blue');

  const loginData = {
    phone: '13800000000',
    password: 'Test123456'
  };

  log(`  请求: POST /api/v1/auth/login`);
  log(`  数据: ${JSON.stringify(loginData)}`);

  try {
    const response = await makeRequest({
      hostname: 'localhost',
      port: 8080,
      path: '/api/v1/auth/login',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    }, loginData);

    if (response.statusCode === 200 && response.body.code === 0) {
      log('✓ 登录接口正常工作', 'green');
      log(`  access_token: ${response.body.data.access_token.substring(0, 50)}...`, 'green');
      log(`  expires_in: ${response.body.data.expires_in}秒`, 'green');
      return response.body.data;
    } else {
      log(`✗ 登录失败: ${JSON.stringify(response.body)}`, 'red');
      return null;
    }
  } catch (error) {
    log(`✗ 请求失败: ${error.message}`, 'red');
    return null;
  }
}

async function testFrontendServer() {
  log('\n========================================', 'blue');
  log('测试3: 前端开发服务器', 'blue');
  log('========================================', 'blue');

  try {
    const response = await makeRequest({
      hostname: 'localhost',
      port: 5175,
      path: '/',
      method: 'GET'
    });

    if (response.statusCode === 200) {
      log('✓ 前端服务正常运行 (http://localhost:5175)', 'green');
      return true;
    } else {
      log(`✗ 前端服务异常: ${response.statusCode}`, 'red');
      return false;
    }
  } catch (error) {
    log(`✗ 无法连接到前端服务: ${error.message}`, 'red');
    log('  提示: 请运行 cd shop-pc && npm run dev', 'yellow');
    return false;
  }
}

async function testCORS() {
  log('\n========================================', 'blue');
  log('测试4: CORS跨域配置', 'blue');
  log('========================================', 'blue');

  try {
    const response = await makeRequest({
      hostname: 'localhost',
      port: 8080,
      path: '/api/v1/auth/login',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Origin': 'http://localhost:5175'
      }
    }, {
      phone: '13800000000',
      password: 'Test123456'
    });

    const corsHeaders = {
      'Access-Control-Allow-Origin': response.headers['access-control-allow-origin'],
      'Access-Control-Allow-Methods': response.headers['access-control-allow-methods'],
      'Access-Control-Allow-Headers': response.headers['access-control-allow-headers']
    };

    if (corsHeaders['Access-Control-Allow-Origin']) {
      log('✓ CORS配置正常', 'green');
      log(`  Headers: ${JSON.stringify(corsHeaders)}`, 'green');
      return true;
    } else {
      log('✗ CORS未配置或配置不正确', 'red');
      log(`  收到的Headers: ${JSON.stringify(corsHeaders)}`, 'yellow');
      return false;
    }
  } catch (error) {
    log(`✗ CORS测试失败: ${error.message}`, 'red');
    return false;
  }
}

async function main() {
  log('\n太积堂系统 - 前后端联调测试', 'blue');
  log('========================================', 'blue');

  const results = {
    backend: await testBackendHealth(),
    frontend: await testFrontendServer(),
    cors: await testCORS(),
    login: null
  };

  if (results.backend && results.frontend) {
    results.login = await testLoginAPI();
  }

  // 总结
  log('\n========================================', 'blue');
  log('测试结果总结', 'blue');
  log('========================================', 'blue');

  const allPassed = results.backend && results.frontend && results.cors && results.login;

  log(`后端服务: ${results.backend ? '✓ 通过' : '✗ 失败'}`, results.backend ? 'green' : 'red');
  log(`前端服务: ${results.frontend ? '✓ 通过' : '✗ 失败'}`, results.frontend ? 'green' : 'red');
  log(`CORS配置: ${results.cors ? '✓ 通过' : '✗ 失败'}`, results.cors ? 'green' : 'red');
  log(`登录接口: ${results.login ? '✓ 通过' : '✗ 失败'}`, results.login ? 'green' : 'red');

  if (allPassed) {
    log('\n✓ 所有测试通过！前后端联调成功', 'green');
    log('\n下一步:', 'blue');
    log('  1. 打开浏览器访问: http://localhost:5175/login', 'blue');
    log('  2. 输入测试账号: 13800000000 / Test123456', 'blue');
    log('  3. 验证登录功能是否正常', 'blue');
    process.exit(0);
  } else {
    log('\n✗ 部分测试失败，请检查上述错误信息', 'red');
    process.exit(1);
  }
}

main().catch(error => {
  log(`\n✗ 测试脚本执行失败: ${error.message}`, 'red');
  console.error(error);
  process.exit(1);
});
