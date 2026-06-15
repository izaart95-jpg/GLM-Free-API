import http from 'node:http';
import { chromium } from 'playwright-core';
import * as cloak from 'cloakbrowser';

const PORT = Number(process.env.PORT || 9876);
const HOST = process.env.HOST || '127.0.0.1';
const CHROMIUM_PATH = process.env.CHROMIUM_PATH || '/usr/bin/chromium';
const SECRET = process.env.SECRET || '';
const ZAI_URL = 'https://chat.z.ai/';
const SCENE_ID = 'didk33e0';

const BROWSER_BACKEND = process.env.BROWSER_BACKEND || 'cloak';

const POOL_SIZE = Number(process.env.POOL_SIZE || 5);
const TOKEN_TTL = Number(process.env.TOKEN_TTL || 240000);
const REFILL_INTERVAL = Number(process.env.REFILL_INTERVAL || 3000);

let browser = null;
let context = null;
let page = null;
let ready = false;
let lastError = '';
let stats = { served: 0, errors: 0, refills: 0 };

const tokenPool = [];
let refilling = false;


async function launchBrowser() {
  console.log(`[provider] Initializing browser backend: ${BROWSER_BACKEND}`);

  if (BROWSER_BACKEND === 'cloak') {
    try {
      const test = await cloak.launch({ headless: true });
      const v = await test.version();
      console.log(`[provider] CloakBrowser ready (chromium ${v})`);
      await test.close();
    } catch (e) {
      console.error('[provider] CloakBrowser init failed:', e.message);
      throw e;
    }
  } else {
    const launchOpts = {
      executablePath: CHROMIUM_PATH,
      headless: true,
      args: [
        '--disable-blink-features=AutomationControlled',
        '--disable-dev-shm-usage',
        '--no-sandbox',
        '--disable-gpu',
        '--disable-extensions',
        '--disable-background-timer-throttling',
        '--disable-renderer-backgrounding',
      ],
    };
    const proxyURL = process.env.PROXY_SERVER || process.env.HTTPS_PROXY || process.env.HTTP_PROXY;
    if (proxyURL) {
      console.log(`[provider] Using proxy: ${proxyURL}`);
      launchOpts.proxy = { server: proxyURL };
    }
    browser = await chromium.launch(launchOpts);

    context = await browser.newContext({
      viewport: { width: 1920, height: 1080 },
      locale: 'zh-CN',
      timezoneId: 'Asia/Shanghai',
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36',
    });

    await context.addInitScript(() => {
      Object.defineProperty(navigator, 'webdriver', { get: () => false });
      Object.defineProperty(navigator, 'plugins', { get: () => [1, 2, 3, 4, 5] });
      Object.defineProperty(navigator, 'languages', { get: () => ['zh-CN', 'zh', 'en-US', 'en'] });
      window.chrome = { runtime: {}, loadTimes: () => ({}) };
      const originalQuery = window.navigator.permissions?.query;
      if (originalQuery) {
        window.navigator.permissions.query = (params) =>
          params.name === 'notifications'
            ? Promise.resolve({ state: Notification.permission })
            : originalQuery(params);
      }
    });
  }

  console.log('[provider] ✓ Browser ready');
  ready = true;
}

async function acquireToken() {
  let localBrowser;

  if (BROWSER_BACKEND === 'cloak') {
    const opts = {
      headless: true,
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36',
    };
    const proxyURL = process.env.PROXY_SERVER || process.env.HTTPS_PROXY || process.env.HTTP_PROXY;
    if (proxyURL) opts.proxy = { server: proxyURL };
    localBrowser = await cloak.launch(opts);
  } else {
    const launchOpts = {
      executablePath: CHROMIUM_PATH,
      headless: true,
      args: [
        '--disable-blink-features=AutomationControlled',
        '--disable-dev-shm-usage',
        '--no-sandbox',
        '--disable-gpu',
        '--disable-extensions',
      ],
    };
    const proxyURL = process.env.PROXY_SERVER || process.env.HTTPS_PROXY || process.env.HTTP_PROXY;
    if (proxyURL) launchOpts.proxy = { server: proxyURL };
    localBrowser = await chromium.launch(launchOpts);
  }

  try {
    const localCtx = await localBrowser.newContext({
      viewport: { width: 1920, height: 1080 },
      locale: 'zh-CN',
      timezoneId: 'Asia/Shanghai',
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36',
    });
    await localCtx.addInitScript(() => {
      Object.defineProperty(navigator, 'webdriver', { get: () => false });
      Object.defineProperty(navigator, 'plugins', { get: () => [1, 2, 3, 4, 5] });
      Object.defineProperty(navigator, 'languages', { get: () => ['zh-CN', 'zh', 'en-US', 'en'] });
      window.chrome = { runtime: {}, loadTimes: () => ({}) };
    });
    const freshPage = await localCtx.newPage();
    await freshPage.goto(ZAI_URL, { waitUntil: 'domcontentloaded', timeout: 30000 });
    await freshPage.waitForTimeout(3000);

    await freshPage.evaluate(async () => {
      if (typeof window.initAliyunCaptcha === 'function') return;
      window.AliyunCaptchaConfig = { region: 'cn', prefix: 'no8xfe' };
      await new Promise((resolve, reject) => {
        const s = document.createElement('script');
        s.src = 'https://o.alicdn.com/captcha-frontend/aliyunCaptcha/AliyunCaptcha.js';
        s.onload = () => resolve();
        s.onerror = () => reject(new Error('SDK load failed'));
        document.head.appendChild(s);
      });
    });
    await freshPage.waitForTimeout(1000);

    const token = await freshPage.evaluate(async (sceneId) => {
      if (typeof window.initAliyunCaptcha !== 'function') {
        throw new Error('initAliyunCaptcha not available after injection');
      }
      const id = 'c-' + Date.now() + '-' + Math.random().toString(36).slice(2, 6);
      const triggerId = 't-' + id;
      const container = document.createElement('div');
      container.id = id;
      container.style.cssText = 'position:fixed;top:-9999px;left:-9999px;width:1px;height:1px;overflow:hidden;';
      document.body.appendChild(container);
      const trigger = document.createElement('button');
      trigger.id = triggerId;
      trigger.style.cssText = 'display:none;';
      container.appendChild(trigger);

      try {
        return await new Promise((resolve, reject) => {
          const timeout = setTimeout(() => reject(new Error('captcha timeout')), 25000);
          window.initAliyunCaptcha({
            SceneId: sceneId,
            mode: 'popup',
            element: `
            button: `
            language: 'cn',
            timeout: 10000,
            delayBeforeSuccess: false,
            success: (token) => { clearTimeout(timeout); resolve(token); },
            fail: () => {},
            onError: (err) => { clearTimeout(timeout); reject(new Error('captcha error: ' + JSON.stringify(err))); },
            onClose: () => {},
          });
          setTimeout(() => trigger.click(), 200);
        });
      } finally {
        try { container.remove(); } catch {}
      }
    }, SCENE_ID);

    return token;
  } finally {
    try { await localBrowser.close(); } catch {}
  }
}


function getValidToken() {
  const now = Date.now();
  while (tokenPool.length > 0 && (now - tokenPool[0].createdAt) > TOKEN_TTL) {
    tokenPool.shift();
  }
  if (tokenPool.length > 0) {
    return tokenPool.shift().token;
  }
  return null;
}

async function refillPool() {
  if (refilling) return;
  refilling = true;
  try {
    const now = Date.now();
    while (tokenPool.length > 0 && (now - tokenPool[0].createdAt) > TOKEN_TTL) {
      tokenPool.shift();
    }
    if (tokenPool.length < POOL_SIZE) {
      try {
        const token = await acquireToken();
        tokenPool.push({ token, createdAt: Date.now() });
        stats.refills++;
      } catch (err) {
        lastError = err.message;
        stats.errors++;
        if (!err.message.includes('captcha timeout')) {
          console.error(`[pool] Refill error: ${err.message}`);
        }
      }
    }
  } finally {
    refilling = false;
  }
}


function sendJson(res, status, data) {
  const body = JSON.stringify(data);
  res.writeHead(status, {
    'Content-Type': 'application/json',
    'Content-Length': Buffer.byteLength(body),
  });
  res.end(body);
}

const server = http.createServer(async (req, res) => {
  if (SECRET && req.headers['x-secret'] !== SECRET) {
    return sendJson(res, 401, { error: 'unauthorized' });
  }

  if (req.method === 'GET' && req.url === '/health') {
    return sendJson(res, 200, {
      ok: ready,
      pool: tokenPool.length,
      stats,
      lastError,
    });
  }

  if (req.method === 'GET' && req.url === '/token') {
    if (!ready) {
      return sendJson(res, 503, { error: 'not ready', lastError });
    }

    const cached = getValidToken();
    if (cached) {
      stats.served++;
      console.log(`[provider] Served cached token (pool: ${tokenPool.length})`);
      return sendJson(res, 200, { ok: true, token: cached, cached: true });
    }

    try {
      const started = Date.now();
      const token = await acquireToken();
      const elapsed = Date.now() - started;
      stats.served++;
      console.log(`[provider] Served fresh token in ${elapsed}ms (pool: ${tokenPool.length})`);
      return sendJson(res, 200, { ok: true, token, cached: false, elapsed_ms: elapsed });
    } catch (err) {
      lastError = err.message;
      stats.errors++;
      console.error(`[provider] Token error: ${err.message}`);
      return sendJson(res, 500, { ok: false, error: err.message });
    }
  }

  sendJson(res, 404, { error: 'Use GET /token or GET /health' });
});


server.listen(PORT, HOST, async () => {
  console.log(`[provider] zai-captcha-provider listening on http:
  console.log(`[provider] Pool size: ${POOL_SIZE}, TTL: ${TOKEN_TTL}ms`);
  try {
    await launchBrowser();
    await refillPool();
    console.log(`[provider] Initial pool filled: ${tokenPool.length} tokens`);
    setInterval(refillPool, REFILL_INTERVAL);
  } catch (err) {
    console.error('[provider] Startup error:', err.message);
    lastError = err.message;
  }
});

process.on('SIGINT', async () => {
  console.log('[provider] Shutting down...');
  if (browser) await browser.close();
  process.exit(0);
});

process.on('SIGTERM', async () => {
  if (browser) await browser.close();
  process.exit(0);
});
