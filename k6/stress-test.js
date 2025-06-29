import http from 'k6/http';
import { check, sleep } from 'k6';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';

// High-concurrency stress test to see ShardStoreGopool benefits
export let options = {
  stages: [
    { duration: '10s', target: 100 },   // warm-up
    { duration: '20s', target: 500 },   // climbing to high load
    { duration: '30s', target: 1000 },  // peak stress test
    { duration: '20s', target: 500 },   // step down
    { duration: '10s', target: 0 },     // cooldown
  ],
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Pre-populate some tasks for read testing
export function setup() {
  console.log('Setting up test data...');
  
  // Create 1000 tasks for read testing
  for (let i = 0; i < 1000; i++) {
    let task = { name: `Setup Task ${i}`, status: i % 2 };
    http.post(`${BASE_URL}/tasks`, JSON.stringify(task), {
      headers: { 'Content-Type': 'application/json' },
    });
  }
  
  console.log('Setup complete - 1000 tasks created');
  return { taskCount: 1000 };
}

export default function (data) {
  // 80% reads, 20% writes (realistic workload)
  const isRead = Math.random() < 0.8;
  
  if (isRead) {
    // Read operations (where our optimization shines)
    let tasks = http.get(`${BASE_URL}/tasks`);
    check(tasks, { 
      'get all tasks success': (r) => r.status === 200,
      'get all tasks fast': (r) => r.timings.duration < 100 // < 100ms
    });
    
    // Test random task ID reads (hot key pattern)
    const taskId = Math.floor(Math.random() * data.taskCount) + 1;
    let task = http.get(`${BASE_URL}/tasks/${taskId}`);
    check(task, { 
      'get task by id success': (r) => r.status === 200 || r.status === 404,
      'get task by id fast': (r) => r.timings.duration < 50 // < 50ms
    });
  } else {
    // Write operations (20%)
    let task = { 
      name: `Stress Task ${Date.now()}`, 
      status: Math.floor(Math.random() * 2) 
    };
    
    let create = http.post(`${BASE_URL}/tasks`, JSON.stringify(task), {
      headers: { 'Content-Type': 'application/json' },
    });
    
    check(create, { 
      'create task success': (r) => r.status === 201,
      'create task fast': (r) => r.timings.duration < 100
    });
  }
  
  // No sleep - maximum throughput test
}

export function handleSummary(data) {
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
  
  return {
    [`/output/stress-${timestamp}.json`]: JSON.stringify(data, null, 2),
    [`/output/stress-${timestamp}.html`]: htmlReport(data),
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
  };
}

function htmlReport(data) {
  const date = new Date().toISOString();
  const duration = data.state.testRunDurationMs / 1000;
  const rps = (data.metrics.http_reqs?.values?.rate || 0).toFixed(2);
  const avgDuration = (data.metrics.http_req_duration?.values?.avg || 0).toFixed(2);
  const p95Duration = (data.metrics.http_req_duration?.values?.['p(95)'] || 0).toFixed(2);
  
  return `<!DOCTYPE html>
<html>
<head>
    <title>K6 Stress Test Report - ShardStoreGopool</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; border-radius: 8px; }
        .metric { margin: 15px 0; padding: 15px; background: #f9f9f9; border-radius: 5px; border-left: 4px solid #667eea; }
        .pass { color: #27ae60; font-weight: bold; }
        .fail { color: #e74c3c; font-weight: bold; }
        .highlight { background: #fff3cd; padding: 10px; border-radius: 5px; margin: 10px 0; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #667eea; color: white; }
        .performance { background: #d4edda; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🚀 K6 Stress Test Report</h1>
        <h2>ShardStoreGopool + ByteDance Optimization</h2>
        <p><strong>Date:</strong> ${date}</p>
        <p><strong>Duration:</strong> ${duration}s</p>
        <p><strong>Peak VUs:</strong> 1000</p>
    </div>
    
    <div class="performance">
        <h2>🎯 Performance Highlights</h2>
        <p><strong>Requests/sec:</strong> ${rps} RPS</p>
        <p><strong>Average Response:</strong> ${avgDuration}ms</p>
        <p><strong>95th Percentile:</strong> ${p95Duration}ms</p>
        <p><strong>Total Requests:</strong> ${data.metrics.http_reqs?.values?.count || 0}</p>
        <p><strong>Failed Requests:</strong> ${((data.metrics.http_req_failed?.values?.rate || 0) * 100).toFixed(2)}%</p>
    </div>
    
    <div class="highlight">
        <strong>Test Pattern:</strong> 80% reads (GetAll + GetByID), 20% writes - simulating realistic workload
        <br><strong>Storage:</strong> ShardStoreGopool with 32 shards + ByteDance gopool optimization
        <br><strong>Target:</strong> High concurrency stress test (up to 1000 VUs)
    </div>
    
    <h2>📊 Detailed Metrics</h2>
    <div class="metric">
        <h3>HTTP Requests</h3>
        <p>Total Requests: ${data.metrics.http_reqs?.values?.count || 0}</p>
        <p>Request Rate: ${rps}/s</p>
        <p>Failed Rate: ${((data.metrics.http_req_failed?.values?.rate || 0) * 100).toFixed(2)}%</p>
    </div>
    
    <div class="metric">
        <h3>Response Times</h3>
        <p>Min: ${(data.metrics.http_req_duration?.values?.min || 0).toFixed(2)}ms</p>
        <p>Average: ${avgDuration}ms</p>
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
    
    <div class="highlight">
        <h3>🔬 Expected Performance with ShardStoreGopool:</h3>
        <ul>
            <li><strong>Read Operations:</strong> ~11.54ns per operation (benchmark level)</li>
            <li><strong>High Concurrency:</strong> Better CPU utilization with per-core pools</li>
            <li><strong>Memory Efficiency:</strong> ShardUnit optimization reduces overhead</li>
            <li><strong>Scalability:</strong> 32 shards optimized for M4 Pro 14-core architecture</li>
        </ul>
    </div>
</body>
</html>`;
}