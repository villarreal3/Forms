<?php
/**
 * Proxy específico para servir imágenes públicas de Forms API.
 * Mapea peticiones desde el frontend a /forms/proxy-forms-images.php/images/{form_id}/{filename}
 * hacia el backend interno http://192.168.2.15:8080/api/images/{form_id}/{filename}.
 */

$BACKEND_BASE = 'http://192.168.2.15:8080/api';

$path = isset($_SERVER['PATH_INFO']) ? $_SERVER['PATH_INFO'] : '';
if ($path === '') {
    http_response_code(400);
    header('Content-Type: application/json');
    echo json_encode(['message' => 'Missing path. Use e.g. /forms/proxy-forms-images.php/images/{form_id}/{filename}']);
    exit;
}

$query = isset($_SERVER['QUERY_STRING']) && $_SERVER['QUERY_STRING'] !== ''
    ? '?' . $_SERVER['QUERY_STRING']
    : '';
$url = $BACKEND_BASE . $path . $query;

$method = $_SERVER['REQUEST_METHOD'] ?? 'GET';

// Para imágenes normalmente solo usamos GET y no enviamos cuerpo,
// pero mantenemos soporte genérico por si se usa en el futuro.
$body = null;
if (in_array($method, ['POST', 'PUT', 'PATCH', 'DELETE'], true)) {
    $body = file_get_contents('php://input');
    if ($body === '') {
        $body = null;
    }
}

$headers = [];
if (isset($_SERVER['HTTP_CONTENT_TYPE']) && $_SERVER['HTTP_CONTENT_TYPE'] !== '') {
    $headers[] = 'Content-Type: ' . $_SERVER['HTTP_CONTENT_TYPE'];
} elseif (!empty($_SERVER['CONTENT_TYPE'])) {
    $headers[] = 'Content-Type: ' . $_SERVER['CONTENT_TYPE'];
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
        'message' => 'No se pudo conectar al servicio de imágenes de formularios.',
        'error'   => curl_strerror($errno),
    ]);
    exit;
}

$parts = preg_split('/\r\n\r\n/', $response, 2);
$responseHeaders = $parts[0] ?? '';
$responseBody = $parts[1] ?? '';

// Reenviar los headers relevantes para archivos binarios
$forwardHeaders = ['Content-Type', 'Content-Length', 'Content-Disposition', 'Cache-Control', 'Last-Modified', 'ETag'];
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

