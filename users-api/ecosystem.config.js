module.exports = {
 apps: [
 {
 name: 'users-api',
 script: './dist/main.js',
 instances: 2, // Número de instancias (puedes usar 'max' para usar todos los CPUs)
 exec_mode: 'cluster',
 env: {
 NODE_ENV: 'production',
 PORT: 3000,
 },
 error_file: './logs/err.log',
 out_file: './logs/out.log',
 log_date_format: 'YYYY-MM-DD HH:mm:ss Z',
 merge_logs: true,
 autorestart: true,
 watch: false,
 max_memory_restart: '1G',
 min_uptime: '10s',
 max_restarts: 10,
 },
 ],
};

