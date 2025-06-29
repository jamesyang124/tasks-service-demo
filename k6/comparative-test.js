import http from 'k6/http';
import { check } from 'k6';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';

// Comparative test across all storage implementations
export let options = {
  stages: [
    { duration: '10s', target: 100 },   // warm-up
    { duration: '20s', target: 300 },   // sustained load
    { duration: '10s', target: 500 },   // peak load
    { duration: '10s', target: 0 },     // cooldown
  ],
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const STORAGE_TYPE = __ENV.STORAGE_TYPE || 'gopool';

export function setup() {
  console.log(`Setting up comparative test for ${STORAGE_TYPE} storage...`);
  
  // Create 1000 tasks for testing
  for (let i = 0; i < 1000; i++) {
    let task = { name: `${STORAGE_TYPE} Task ${i}`, status: i % 2 };
    http.post(`${BASE_URL}/tasks`, JSON.stringify(task), {
      headers: { 'Content-Type': 'application/json' },
    });
  }
  
  console.log(`Setup complete - ${STORAGE_TYPE} storage ready with 1000 tasks`);
  return { taskCount: 1000, storageType: STORAGE_TYPE };
}

export default function (data) {
  // 80% reads, 20% writes (realistic workload)
  const isRead = Math.random() < 0.8;
  
  if (isRead) {
    // Read operations (where storage differences are most visible)
    if (Math.random() < 0.7) {
      // 70% - Individual task reads (hot key pattern)
      const taskId = Math.floor(Math.random() * data.taskCount) + 1;
      let response = http.get(`${BASE_URL}/tasks/${taskId}`);
      check(response, { 
        [`${data.storageType} read success`]: (r) => r.status === 200,
        [`${data.storageType} read fast`]: (r) => r.timings.duration < 100
      });
    } else {
      // 10% - Get all tasks (bulk operation test)
      let response = http.get(`${BASE_URL}/tasks`);
      check(response, { 
        [`${data.storageType} getall success`]: (r) => r.status === 200,
        [`${data.storageType} getall reasonable`]: (r) => r.timings.duration < 500
      });
    }
  } else {
    // 20% writes
    let task = { name: `${data.storageType} Load Task ${Date.now()}`, status: Math.floor(Math.random() * 2) };
    let response = http.post(`${BASE_URL}/tasks`, JSON.stringify(task), {
      headers: { 'Content-Type': 'application/json' },
    });
    check(response, { 
      [`${data.storageType} write success`]: (r) => r.status === 201,
      [`${data.storageType} write fast`]: (r) => r.timings.duration < 200
    });
  }
}

export function handleSummary(data) {
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
  const storageType = __ENV.STORAGE_TYPE || 'gopool';
  
  const rps = (data.metrics.http_reqs?.values?.rate || 0).toFixed(2);
  const avgDuration = (data.metrics.http_req_duration?.values?.avg || 0).toFixed(2);
  const p95Duration = (data.metrics.http_req_duration?.values?.['p(95)'] || 0).toFixed(2);
  const errorRate = ((data.metrics.http_req_failed?.values?.rate || 0) * 100).toFixed(2);
  
  console.log(`\n🔬 ${storageType.toUpperCase()} STORAGE RESULTS:`);
  console.log(`RPS: ${rps}`);
  console.log(`Avg Response: ${avgDuration}ms`);
  console.log(`95th Percentile: ${p95Duration}ms`);
  console.log(`Error Rate: ${errorRate}%`);
  console.log(`Total Requests: ${data.metrics.http_reqs?.values?.count || 0}`);
  
  return {
    [`/output/comparative-${storageType}-${timestamp}.json`]: JSON.stringify(data, null, 2),
    [`/output/comparative-${storageType}-${timestamp}.html`]: htmlReport(data, storageType),
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
  };
}

function htmlReport(data, storageType) {
  const date = new Date().toISOString();
  const duration = data.state.testRunDurationMs / 1000;
  const rps = (data.metrics.http_reqs?.values?.rate || 0).toFixed(2);
  const avgDuration = (data.metrics.http_req_duration?.values?.avg || 0).toFixed(2);
  const p95Duration = (data.metrics.http_req_duration?.values?.['p(95)'] || 0).toFixed(2);
  const errorRate = ((data.metrics.http_req_failed?.values?.rate || 0) * 100).toFixed(2);
  
  // Storage-specific styling and descriptions
  const storageInfo = {
    'gopool': {
      color: '#4CAF50',
      name: 'ShardStoreGopool',
      description: 'ByteDance gopool + ShardUnit optimization (11.54ns reads)',
      expected: 'Highest performance with per-core worker pools'
    },
    'shard': {
      color: '#2196F3', 
      name: 'ShardStore',
      description: 'Dedicated workers + ShardUnit (12.42ns reads)',
      expected: 'High performance with dedicated worker pattern'
    },
    'memory': {
      color: '#FF9800',
      name: 'MemoryStore', 
      description: 'Single mutex implementation (130ns reads)',
      expected: 'Baseline performance, limited concurrency'
    },
    'bigcache': {
      color: '#9C27B0',
      name: 'BigCacheStore',
      description: 'Off-heap cache implementation (65ns reads)',
      expected: 'Medium performance with GC advantages'
    }
  };
  
  const info = storageInfo[storageType] || storageInfo['gopool'];
  
  return `<!DOCTYPE html>
<html>
<head>
    <title>${info.name} Performance Test Results</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: linear-gradient(135deg, ${info.color} 0%, ${info.color}AA 100%); color: white; padding: 20px; border-radius: 8px; }
        .metric { margin: 15px 0; padding: 15px; background: #f9f9f9; border-radius: 5px; border-left: 4px solid ${info.color}; }
        .pass { color: #4CAF50; font-weight: bold; }
        .fail { color: #F44336; font-weight: bold; }
        .highlight { background: #f5f5f5; padding: 15px; border-radius: 5px; margin: 15px 0; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: ${info.color}; color: white; }
        .performance { background: #e8f5e8; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid ${info.color}; }
    </style>
</head>
<body>
    <div class="header">
        <h1>📊 ${info.name} Performance Test</h1>
        <h3>${info.description}</h3>
        <p><strong>Date:</strong> ${date}</p>
        <p><strong>Duration:</strong> ${duration}s</p>
        <p><strong>Peak Load:</strong> 500 VUs</p>
    </div>
    
    <div class="performance">
        <h2>🎯 Key Performance Metrics</h2>
        <p><strong>Requests/sec:</strong> ${rps} RPS</p>
        <p><strong>Average Response:</strong> ${avgDuration}ms</p>
        <p><strong>95th Percentile:</strong> ${p95Duration}ms</p>
        <p><strong>Error Rate:</strong> ${errorRate}%</p>
        <p><strong>Total Requests:</strong> ${data.metrics.http_reqs?.values?.count || 0}</p>
    </div>
    
    <div class="highlight">
        <strong>Test Pattern:</strong> 80% reads (70% GetByID + 10% GetAll), 20% writes
        <br><strong>Dataset:</strong> 1,000 tasks with realistic workload simulation
        <br><strong>Expected:</strong> ${info.expected}
    </div>
    
    <h2>📈 Detailed Response Times</h2>
    <div class="metric">
        <h3>HTTP Request Duration</h3>
        <p>Min: ${(data.metrics.http_req_duration?.values?.min || 0).toFixed(2)}ms</p>
        <p>Average: ${avgDuration}ms</p>
        <p>Median (p50): ${(data.metrics.http_req_duration?.values?.med || 0).toFixed(2)}ms</p>
        <p>90th percentile: ${(data.metrics.http_req_duration?.values?.['p(90)'] || 0).toFixed(2)}ms</p>
        <p>95th percentile: ${p95Duration}ms</p>
        <p>99th percentile: ${(data.metrics.http_req_duration?.values?.['p(99)'] || 0).toFixed(2)}ms</p>
        <p>Max: ${(data.metrics.http_req_duration?.values?.max || 0).toFixed(2)}ms</p>
    </div>
    
    <h2>✅ Check Results</h2>
    <table>
        <tr><th>Check</th><th>Passes</th><th>Failures</th><th>Rate</th></tr>
        ${Object.entries(data.metrics)
          .filter(([key]) => key.startsWith('checks'))
          .map(([key, metric]) => {
            const name = key.replace('checks{', '').replace('}', '').replace('name:', '');
            const rate = (metric.values.rate * 100).toFixed(2);
            const status = metric.values.rate >= 0.95 ? 'pass' : 'fail';
            return '<tr class="' + status + '">' +
              '<td>' + name + '</td>' +
              '<td>' + metric.values.passes + '</td>' +
              '<td>' + metric.values.fails + '</td>' +
              '<td>' + rate + '%</td>' +
            '</tr>';
          }).join('')}
    </table>
</body>
</html>`;
}