import http from "k6/http";
import { check, group, sleep } from "k6";

const data = JSON.parse(open("../data/test-data.json"));

const baseUrl = (__ENV.TARGET_BASE_URL || "http://localhost:8080").replace(/\/$/, "");
const profile = __ENV.PROFILE || "normal";
const scenarioName = __ENV.SCENARIO || "full";
const allowWrites = (__ENV.ALLOW_WRITES || "false").toLowerCase() === "true";
const allowProduction = (__ENV.ALLOW_PRODUCTION_LOAD_TEST || "false").toLowerCase() === "true";

const authToken = __ENV.AUTH_TOKEN || "";
const apiPublicKey = __ENV.API_PUBLIC_KEY || "";
const apiPrivateKey = __ENV.API_PRIVATE_KEY || "";
const venueId = __ENV.VENUE_ID || data.ids.venue_id;
const locationId = __ENV.LOCATION_ID || data.ids.location_id;
const productId = __ENV.PRODUCT_ID || data.ids.product_id;
const surveyId = __ENV.SURVEY_ID || data.ids.survey_id;

const profiles = {
  normal: {
    scenarios: {
      public_read: { executor: "constant-vus", vus: 5, duration: "2m", exec: "publicReadFlow" },
      app_read: { executor: "constant-vus", vus: 10, duration: "2m", exec: "appReadFlow" },
      admin_read: { executor: "constant-vus", vus: 3, duration: "2m", exec: "adminReadFlow" },
    },
  },
  peak: {
    scenarios: {
      public_read: { executor: "ramping-vus", stages: [{ duration: "1m", target: 20 }, { duration: "3m", target: 20 }, { duration: "1m", target: 0 }], exec: "publicReadFlow" },
      app_read: { executor: "ramping-vus", stages: [{ duration: "1m", target: 50 }, { duration: "3m", target: 50 }, { duration: "1m", target: 0 }], exec: "appReadFlow" },
      admin_read: { executor: "ramping-vus", stages: [{ duration: "1m", target: 10 }, { duration: "3m", target: 10 }, { duration: "1m", target: 0 }], exec: "adminReadFlow" },
    },
  },
  burst: {
    scenarios: {
      burst_public: { executor: "ramping-arrival-rate", startRate: 10, timeUnit: "1s", preAllocatedVUs: 60, maxVUs: 120, stages: [{ duration: "30s", target: 150 }, { duration: "30s", target: 150 }, { duration: "1m", target: 20 }], exec: "publicReadFlow" },
      burst_app: { executor: "ramping-arrival-rate", startRate: 10, timeUnit: "1s", preAllocatedVUs: 80, maxVUs: 160, stages: [{ duration: "30s", target: 200 }, { duration: "30s", target: 200 }, { duration: "1m", target: 30 }], exec: "appReadFlow" },
    },
  },
  error_retry: {
    scenarios: {
      error_retry: { executor: "constant-vus", vus: 5, duration: "2m", exec: "errorRetryFlow" },
    },
    thresholds: {
      http_req_failed: ["rate<0.95"],
      http_req_duration: ["p(95)<1500", "p(99)<2500"],
    },
  },
};

export const options = {
  ...(profiles[profile] || profiles.normal),
  thresholds: profiles[profile] && profiles[profile].thresholds ? profiles[profile].thresholds : {
    http_req_failed: ["rate<0.05"],
    http_req_duration: ["avg<500", "p(95)<1000", "p(99)<2000"],
    "http_req_duration{flow:public}": ["p(95)<800"],
    "http_req_duration{flow:app}": ["p(95)<1000"],
    "http_req_duration{flow:admin}": ["p(95)<1200"],
  },
  summaryTrendStats: ["avg", "min", "med", "p(90)", "p(95)", "p(99)", "max"],
};

export function setup() {
  guardEnvironment();
  return { baseUrl };
}

export function publicReadFlow() {
  if (!shouldRun("public")) return;

  group("public read flow", () => {
    expectOK("GET /health", request("GET", "/health", null, { flow: "public", endpoint: "GET /health" }));
    expectOK("GET /version", request("GET", "/version", null, { flow: "public", endpoint: "GET /version" }));
    expectOK("GET /public/v1/venue-info", request("GET", `/public/v1/venue-info?public_key=${encodeURIComponent(apiPublicKey)}`, null, { flow: "public", endpoint: "GET /public/v1/venue-info" }, [200, 400, 404]));
    expectOK("GET /public/v1/visitor-surveys", request("GET", `/public/v1/visitor-surveys?public_key=${encodeURIComponent(apiPublicKey)}`, null, { flow: "public", endpoint: "GET /public/v1/visitor-surveys" }, [200, 400, 404]));
    expectOK("POST /api/v1/public/venues/:id/events", request("POST", `/api/v1/public/venues/${venueId}/events`, data.payloads.track_event_valid, { flow: "public", endpoint: "POST /api/v1/public/venues/:id/events" }, [200, 404]));
    expectOK("POST /api/v1/public/venues/:id/searches", request("POST", `/api/v1/public/venues/${venueId}/searches`, data.payloads.track_search_valid, { flow: "public", endpoint: "POST /api/v1/public/venues/:id/searches" }, [200, 404]));
  });

  sleep(1);
}

export function appReadFlow() {
  if (!shouldRun("app")) return;

  group("app read flow", () => {
    const headers = apiHeaders();
    expectOK("GET /app/v1/locations", request("GET", "/app/v1/locations?page=1&page_size=20", null, { flow: "app", endpoint: "GET /app/v1/locations" }, [200, 401], headers));
    expectOK("GET /app/v1/products", request("GET", "/app/v1/products?page=1&page_size=20", null, { flow: "app", endpoint: "GET /app/v1/products" }, [200, 401], headers));
    expectOK("GET /app/v1/events", request("GET", "/app/v1/events?page=1&page_size=20", null, { flow: "app", endpoint: "GET /app/v1/events" }, [200, 401], headers));
    expectOK("GET /app/v1/articles", request("GET", "/app/v1/articles?page=1&page_size=20", null, { flow: "app", endpoint: "GET /app/v1/articles" }, [200, 401], headers));
    expectOK("GET /app/v1/search-options", request("GET", "/app/v1/search-options?q=coffee", null, { flow: "app", endpoint: "GET /app/v1/search-options" }, [200, 401], headers));
    expectOK("GET /app/v1/promotions", request("GET", "/app/v1/promotions?page=1&page_size=20", null, { flow: "app", endpoint: "GET /app/v1/promotions" }, [200, 401], headers));

    if (locationId) {
      expectOK("GET /app/v1/locations/:id", request("GET", `/app/v1/locations/${locationId}`, null, { flow: "app", endpoint: "GET /app/v1/locations/:id" }, [200, 401, 404], headers));
    }
    if (productId) {
      expectOK("GET /app/v1/products/:id", request("GET", `/app/v1/products/${productId}`, null, { flow: "app", endpoint: "GET /app/v1/products/:id" }, [200, 401, 404], headers));
    }
  });

  sleep(1);
}

export function adminReadFlow() {
  if (!shouldRun("admin")) return;

  group("admin read flow", () => {
    const headers = bearerHeaders();
    expectOK("GET /api/v1/profile", request("GET", "/api/v1/profile", null, { flow: "admin", endpoint: "GET /api/v1/profile" }, [200, 401], headers));
    expectOK("GET /api/v1/venues", request("GET", "/api/v1/venues?page=1&page_size=20", null, { flow: "admin", endpoint: "GET /api/v1/venues" }, [200, 401], headers));
    expectOK("GET /api/v1/venues/:id", request("GET", `/api/v1/venues/${venueId}`, null, { flow: "admin", endpoint: "GET /api/v1/venues/:id" }, [200, 401, 403, 404], headers));
    expectOK("GET /api/v1/venues/:id/locations", request("GET", `/api/v1/venues/${venueId}/locations?page=1&page_size=20`, null, { flow: "admin", endpoint: "GET /api/v1/venues/:id/locations" }, [200, 401, 403, 404], headers));
    expectOK("GET /api/v1/venues/:id/products", request("GET", `/api/v1/venues/${venueId}/products?page=1&page_size=20`, null, { flow: "admin", endpoint: "GET /api/v1/venues/:id/products" }, [200, 401, 403, 404], headers));
    expectOK("GET /api/v1/venues/:id/analytics/searches", request("GET", `/api/v1/venues/${venueId}/analytics/searches?page=1&page_size=20`, null, { flow: "admin", endpoint: "GET /api/v1/venues/:id/analytics/searches" }, [200, 401, 403, 404], headers));
  });

  if (allowWrites) {
    group("admin write smoke flow", () => {
      const headers = bearerHeaders();
      expectOK("POST /api/v1/venues/:id/categories", request("POST", `/api/v1/venues/${venueId}/categories`, data.payloads.category_valid, { flow: "admin", endpoint: "POST /api/v1/venues/:id/categories" }, [200, 201, 401, 403, 404], headers));
      expectOK("POST /api/v1/venues/:id/locations", request("POST", `/api/v1/venues/${venueId}/locations`, data.payloads.location_valid, { flow: "admin", endpoint: "POST /api/v1/venues/:id/locations" }, [200, 201, 401, 403, 404], headers));
    });
  }

  sleep(1);
}

export function errorRetryFlow() {
  if (!shouldRun("error_retry")) return;

  group("error and retry flow", () => {
    expectFailure("GET /api/v1/profile missing bearer", request("GET", "/api/v1/profile", null, { flow: "error", endpoint: "GET /api/v1/profile missing bearer" }, [401]));
    expectFailure("GET /app/v1/locations invalid api key", request("GET", "/app/v1/locations", null, { flow: "error", endpoint: "GET /app/v1/locations invalid api key" }, [401], { Authorization: "ApiKey invalid:invalid" }));
    expectFailure("POST /api/v1/auth/login invalid body", request("POST", "/api/v1/auth/login", data.payloads.login_invalid, { flow: "error", endpoint: "POST /api/v1/auth/login invalid body" }, [400, 401]));
    expectFailure("GET /public/v1/venues/:id/information invalid uuid", request("GET", "/public/v1/venues/not-a-uuid/information", null, { flow: "error", endpoint: "GET /public/v1/venues/:id/information invalid uuid" }, [400]));
    retryableRequest("GET /health retry wrapper", "/health", { flow: "error", endpoint: "GET /health retry wrapper" });
  });

  sleep(1);
}

function request(method, path, body, tags, expectedStatuses = [200], headers = {}) {
  const params = {
    headers: { "Content-Type": "application/json", ...headers },
    tags,
  };
  const payload = body === null || body === undefined ? null : JSON.stringify(body);
  const res = http.request(method, `${baseUrl}${path}`, payload, params);
  res.expectedStatuses = expectedStatuses;
  return res;
}

function retryableRequest(name, path, tags) {
  let res;
  for (let attempt = 1; attempt <= 3; attempt += 1) {
    res = request("GET", path, null, { ...tags, attempt: String(attempt) }, [200, 429, 500, 502, 503, 504]);
    if (![429, 500, 502, 503, 504].includes(res.status)) break;
    sleep(attempt);
  }
  expectOK(name, res);
}

function expectOK(name, res) {
  return check(res, {
    [`${name}: expected status`]: (r) => r.expectedStatuses.includes(r.status),
    [`${name}: response envelope valid when JSON`]: (r) => isJSON(r) ? hasEnvelope(r) : true,
  });
}

function expectFailure(name, res) {
  return check(res, {
    [`${name}: expected failure status`]: (r) => r.expectedStatuses.includes(r.status),
    [`${name}: error envelope`]: (r) => isJSON(r) ? hasEnvelope(r) && json(r).code !== 0 : true,
  });
}

function bearerHeaders() {
  return authToken ? { Authorization: `Bearer ${authToken}` } : {};
}

function apiHeaders() {
  return apiPublicKey && apiPrivateKey ? { Authorization: `ApiKey ${apiPublicKey}:${apiPrivateKey}` } : {};
}

function shouldRun(name) {
  return scenarioName === "full" || scenarioName === name;
}

function guardEnvironment() {
  const lower = baseUrl.toLowerCase();
  const looksProduction = lower.includes("prod") || lower.includes("production") || lower.includes("api.digimap");
  if (looksProduction && !allowProduction) {
    throw new Error("Refusing to run against a production-looking target. Set ALLOW_PRODUCTION_LOAD_TEST=true to override.");
  }
  if ((scenarioName === "full" || scenarioName === "app") && (!apiPublicKey || !apiPrivateKey)) {
    console.warn("API_PUBLIC_KEY/API_PRIVATE_KEY are not set; app API requests should return 401.");
  }
  if ((scenarioName === "full" || scenarioName === "admin") && !authToken) {
    console.warn("AUTH_TOKEN is not set; admin API requests should return 401.");
  }
}

function isJSON(res) {
  return (res.headers["Content-Type"] || "").includes("application/json");
}

function hasEnvelope(res) {
  const body = json(res);
  return body && Object.prototype.hasOwnProperty.call(body, "code");
}

function json(res) {
  try {
    return res.json();
  } catch (_) {
    return null;
  }
}

export function handleSummary(summary) {
  const prefix = `${profile}-${scenarioName}`;
  return {
    [`reports/${prefix}-summary.json`]: JSON.stringify(summary, null, 2),
    [`reports/${prefix}-summary.html`]: htmlSummary(summary),
    stdout: textSummary(summary),
  };
}

function textSummary(summary) {
  const metrics = summary.metrics;
  const lines = [
    "",
    "Load test summary",
    `Profile: ${profile}`,
    `Scenario: ${scenarioName}`,
    `Requests: ${metricCount(metrics.http_reqs)}`,
    `Request rate: ${metricRate(metrics.http_reqs)} req/s`,
    `Average latency: ${metricValue(metrics.http_req_duration, "avg")} ms`,
    `P95 latency: ${metricValue(metrics.http_req_duration, "p(95)")} ms`,
    `P99 latency: ${metricValue(metrics.http_req_duration, "p(99)")} ms`,
    `Error rate: ${metricRate(metrics.http_req_failed)}`,
    `Failed requests: ${metricPasses(metrics.http_req_failed)}`,
    `Reports: load-tests/reports/${profile}-${scenarioName}-summary.json, load-tests/reports/${profile}-${scenarioName}-summary.html`,
    "",
  ];
  return lines.join("\n");
}

function metricValue(metric, key) {
  return metric && metric.values && metric.values[key] !== undefined ? metric.values[key].toFixed(2) : "n/a";
}

function metricCount(metric) {
  return metric && metric.values && metric.values.count !== undefined ? metric.values.count : "n/a";
}

function metricRate(metric) {
  return metric && metric.values && metric.values.rate !== undefined ? metric.values.rate.toFixed(2) : "n/a";
}

function metricPasses(metric) {
  return metric && metric.values && metric.values.passes !== undefined ? metric.values.passes : "n/a";
}

function htmlSummary(summary) {
  const metrics = summary.metrics;
  const rows = [
    ["Requests", metricCount(metrics.http_reqs)],
    ["Request rate", `${metricRate(metrics.http_reqs)} req/s`],
    ["Average latency", `${metricValue(metrics.http_req_duration, "avg")} ms`],
    ["P95 latency", `${metricValue(metrics.http_req_duration, "p(95)")} ms`],
    ["P99 latency", `${metricValue(metrics.http_req_duration, "p(99)")} ms`],
    ["Error rate", metricRate(metrics.http_req_failed)],
    ["Failed requests", metricPasses(metrics.http_req_failed)],
  ].map(([name, value]) => `<tr><th>${name}</th><td>${value}</td></tr>`).join("\n");

  return `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Digimap Load Test Summary</title>
  <style>
    body { font-family: sans-serif; margin: 32px; color: #17202a; }
    table { border-collapse: collapse; min-width: 420px; }
    th, td { border: 1px solid #d6dbdf; padding: 10px 12px; text-align: left; }
    th { background: #f4f6f7; }
    code { background: #f4f6f7; padding: 2px 4px; }
  </style>
</head>
<body>
  <h1>Digimap Load Test Summary</h1>
  <p><strong>Profile:</strong> <code>${profile}</code> <strong>Scenario:</strong> <code>${scenarioName}</code></p>
  <table>${rows}</table>
</body>
</html>`;
}
