import http from 'k6/http';
import { check } from 'k6';

// Quick load test to find the limits
export let options = {
  stages: [
    { duration: '10s', target: 500 },   // ramp up to 500
    { duration: '10s', target: 1000 },  // ramp up to 1000
    { duration: '10s', target: 2000 },  // ramp up to 2000
    { duration: '10s', target: 3000 },  // push to 3000
    { duration: '10s', target: 0 },     // cool down
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'],  // 95% under 100ms
    http_req_failed: ['rate<0.1'],     // <10% errors
  },
};

const BASE_URL = 'http://localhost:8080';

export function setup() {
  console.log('Setting up 1000 tasks for load testing...');
  
  // Create 1000 tasks for read testing
  for (let i = 0; i < 1000; i++) {
    let task = { name: `Load Test Task ${i}`, status: i % 2 };
    http.post(`${BASE_URL}/tasks`, JSON.stringify(task), {
      headers: { 'Content-Type': 'application/json' },
    });
  }
  
  return { taskCount: 1000 };
}

export default function (data) {
  // 90% reads, 10% writes - read optimization focus
  if (Math.random() < 0.9) {
    // Read operations
    if (Math.random() < 0.7) {
      // 70% - Individual task reads (hot key pattern)
      const taskId = Math.floor(Math.random() * data.taskCount) + 1;
      let response = http.get(`${BASE_URL}/tasks/${taskId}`);
      check(response, { 
        'read task success': (r) => r.status === 200,
        'read task fast': (r) => r.timings.duration < 50
      });
    } else {
      // 20% - Get all tasks
      let response = http.get(`${BASE_URL}/tasks`);
      check(response, { 
        'get all success': (r) => r.status === 200,
        'get all reasonable': (r) => r.timings.duration < 200
      });
    }
  } else {
    // 10% writes
    let task = { name: `Load Task ${Date.now()}`, status: Math.floor(Math.random() * 2) };
    let response = http.post(`${BASE_URL}/tasks`, JSON.stringify(task), {
      headers: { 'Content-Type': 'application/json' },
    });
    check(response, { 
      'create success': (r) => r.status === 201,
      'create fast': (r) => r.timings.duration < 100
    });
  }
}

export function handleSummary(data) {
  console.log(`\n🚀 LOAD TEST RESULTS:`);
  console.log(`Peak RPS: ${(data.metrics.http_reqs?.values?.rate || 0).toFixed(0)}`);
  console.log(`Avg Response: ${(data.metrics.http_req_duration?.values?.avg || 0).toFixed(2)}ms`);
  console.log(`95th Percentile: ${(data.metrics.http_req_duration?.values?.['p(95)'] || 0).toFixed(2)}ms`);
  console.log(`Error Rate: ${((data.metrics.http_req_failed?.values?.rate || 0) * 100).toFixed(2)}%`);
  console.log(`Total Requests: ${data.metrics.http_reqs?.values?.count || 0}`);
  
  return {
    'stdout': JSON.stringify(data, null, 2),
  };
}