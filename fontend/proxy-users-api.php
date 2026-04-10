<?php
/**
 * Proxy para Users API (192.168.2.250:3001).
 * Recibe peticiones del frontend (mismo origen HTTPS) y las reenvía al backend interno.
 * No requiere cambios en la configuración de Apache.
 */

$BACKEND_BASE = 'http://192.168.2.250:3001/api/v1';

$path = isset($_SERVER['PATH_INFO']) ? $_SERVER['PATH_INFO'] : '';
if ($path === '') {
    http_response_code(400);
    header('Content-Type: application/json');
    echo json_encode(['message' => 'Missing path. Use e.g. /crece/proxy-users-api.php/auth/login']);
    exit;
}

$query = isset($_SERVER['QUERY_STRING']) && $_SERVER['QUERY_STRING'] !== ''
    ? '?' . $_SERVER['QUERY_STRING']
    : '';
$url = $BACKEND_BASE . $path . $query;

$method = $_SERVER['REQUEST_METHOD'] ?? 'GET';

$body = null;
if (in_array($method, ['POST', 'PUT', 'PATCH', 'DELETE'], true)) {
    $body = file_get_contents('php://input');
    if ($body === '') {
        $body = null;
    }
}

// Authorization: Apache a veces no pasa este header a PHP; probar varias fuentes.
$authHeader = null;
if (!empty($_SERVER['HTTP_AUTHORIZATION'])) {
    $authHeader = $_SERVER['HTTP_AUTHORIZATION'];
} elseif (!empty($_SERVER['REDIRECT_HTTP_AUTHORIZATION'])) {
    $authHeader = $_SERVER['REDIRECT_HTTP_AUTHORIZATION'];
} elseif (function_exists('getallheaders')) {
    $all = getallheaders();
    foreach ($all as $name => $value) {
        if (strcasecmp($name, 'Authorization') === 0) {
            $authHeader = $value;
            break;
        }
    }
}

$headers = [];
if (isset($_SERVER['HTTP_CONTENT_TYPE']) && $_SERVER['HTTP_CONTENT_TYPE'] !== '') {
    $headers[] = 'Content-Type: ' . $_SERVER['HTTP_CONTENT_TYPE'];
} elseif ($body !== null) {
    $headers[] = 'Content-Type: application/json';
}
if ($authHeader !== null && $authHeader !== '') {
    $headers[] = 'Authorization: ' . $authHeader;
}

$ch = curl_init($url);
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_HEADER, true);
curl_setopt($ch, CURLOPT_TIMEOUT, 60);
curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
if ($body !== null) {
    curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
}

$response = curl_exec($ch);
$errno = curl_errno($ch);
$httpCode = (int) curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($errno) {
    http_response_code(502);
    header('Content-Type: application/json');
    echo json_encode([
        'message' => 'No se pudo conectar al servicio de usuarios.',
        'error'   => curl_strerror($errno),
    ]);
    exit;
}

$parts = preg_split('/\r\n\r\n/', $response, 2);
$responseHeaders = $parts[0] ?? '';
$responseBody = $parts[1] ?? '';

$forwardHeaders = ['Content-Type', 'Content-Disposition'];
foreach (explode("\r\n", $responseHeaders) as $line) {
    if (stripos($line, 'HTTP/') === 0) {
        continue;
    }
    $colon = strpos($line, ':');
    if ($colon === false) {
        continue;
    }
    $name = trim(substr($line, 0, $colon));
    foreach ($forwardHeaders as $allow) {
        if (strcasecmp($name, $allow) === 0) {
            header($line);
            break;
        }
    }
}

http_response_code($httpCode);
echo $responseBody;
