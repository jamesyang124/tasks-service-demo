import http from 'k6/http';
import { check, sleep } from 'k6';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';


// load testing mode
// export let options = {
//   vus: 10, // 10 virtual users
//   duration: '30s', // run for 30 seconds
// };

// stress
export let options = {
  stages: [
    { duration: '10s', target: 10 },   // warm-up
    { duration: '10s', target: 50 },   // climbing
    { duration: '10s', target: 100 },  // peak test
    { duration: '10s', target: 0 },    // cooldown
  ],
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Test all endpoints
  let tasks = http.get(`${BASE_URL}/tasks`);
  check(tasks, { 'get tasks': (r) => r.status === 200 });

  let task = { name: 'Test task', status: 0 };
  let create = http.post(`${BASE_URL}/tasks`, JSON.stringify(task), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  if (check(create, { 'create task': (r) => r.status === 201 })) {
    let createdTask = JSON.parse(create.body);
    
    // Update task
    let update = http.put(`${BASE_URL}/tasks/${createdTask.id}`, 
      JSON.stringify({ name: 'Updated task', status: 1 }), {
      headers: { 'Content-Type': 'application/json' },
    });
    check(update, { 'update task': (r) => r.status === 200 });
    
    // Delete task
    let del = http.del(`${BASE_URL}/tasks/${createdTask.id}`);
    check(del, { 'delete task': (r) => r.status === 204 });
  }
  
  sleep(1);
}

export function handleSummary(data) {
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
  
  return {
    '/output/k6-summary.json': JSON.stringify(data, null, 2),
    '/output/k6-report.html': htmlReport(data),
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
  };
}

function htmlReport(data) {
  const date = new Date().toISOString();
  const duration = data.state.testRunDurationMs / 1000;
  
  return `<!DOCTYPE html>
<html>
<head>
    <title>K6 Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f5f5f5; padding: 20px; border-radius: 5px; }
        .metric { margin: 10px 0; padding: 10px; background: #f9f9f9; border-radius: 3px; }
        .pass { color: green; }
        .fail { color: red; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>K6 Load Test Report</h1>
        <p><strong>Date:</strong> ${date}</p>
        <p><strong>Duration:</strong> ${duration}s</p>
        <p><strong>VUs:</strong> ${data.options.scenarios?.default?.vus || 'N/A'}</p>
    </div>
    
    <h2>Test Results</h2>
    <div class="metric">
        <h3>HTTP Requests</h3>
        <p>Total: ${data.metrics.http_reqs?.values?.count || 0}</p>
        <p>Rate: ${(data.metrics.http_reqs?.values?.rate || 0).toFixed(2)}/s</p>
        <p>Failed: ${(data.metrics.http_req_failed?.values?.rate * 100 || 0).toFixed(2)}%</p>
    </div>
    
    <div class="metric">
        <h3>Response Times</h3>
        <p>Average: ${(data.metrics.http_req_duration?.values?.avg || 0).toFixed(2)}ms</p>
        <p>95th percentile: ${(data.metrics.http_req_duration?.values?.['p(95)'] || 0).toFixed(2)}ms</p>
        <p>Max: ${(data.metrics.http_req_duration?.values?.max || 0).toFixed(2)}ms</p>
    </div>
    
    <h2>Check Results</h2>
    <table>
        <tr><th>Check</th><th>Passes</th><th>Failures</th><th>Rate</th></tr>
        ${Object.entries(data.metrics)
          .filter(([key]) => key.startsWith('checks'))
          .map(([key, metric]) => {
            const name = key.replace('checks{', '').replace('}', '');
            const rate = (metric.values.rate * 100).toFixed(2);
            const status = metric.values.rate === 1 ? 'pass' : 'fail';
            return '<tr class="' + status + '">' +
              '<td>' + name + '</td>' +
              '<td>' + metric.values.passes + '</td>' +
              '<td>' + metric.values.fails + '</td>' +
              '<td>' + rate + '%</td>' +
            '</tr>';
          }).join('')}
    </table>
    
    <h2>Raw Data</h2>
    <pre>${JSON.stringify(data, null, 2)}</pre>
</body>
</html>`;
}