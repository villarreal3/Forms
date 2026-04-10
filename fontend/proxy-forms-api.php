<?php
/**
 * Proxy para Forms API (http://192.168.2.15:8080).
 * Recibe peticiones del frontend (mismo origen HTTPS) y las reenvía al backend interno.
 * No requiere cambios en la configuración de Apache.
 * Incluye endpoints públicos (ej. /public/forms, /forms/submit) y privados (/admin/*).
 * Endpoints privados: el frontend envía Authorization: Bearer <token> (token del login en Users API); este proxy lo reenvía.
 */

$BACKEND_BASE = 'http://192.168.2.15:8080/api';

$path = isset($_SERVER['PATH_INFO']) ? $_SERVER['PATH_INFO'] : '';
if ($path === '') {
    http_response_code(400);
    header('Content-Type: application/json');
    echo json_encode(['message' => 'Missing path. Use e.g. /crece/proxy-forms-api.php/public/forms']);
    exit;
}

$query = isset($_SERVER['QUERY_STRING']) && $_SERVER['QUERY_STRING'] !== ''
    ? '?' . $_SERVER['QUERY_STRING']
    : '';
$url = $BACKEND_BASE . $path . $query;

$method = $_SERVER['REQUEST_METHOD'] ?? 'GET';

// Content-Type: para detectar multipart (php://input no está disponible en multipart).
$contentType = isset($_SERVER['HTTP_CONTENT_TYPE']) ? $_SERVER['HTTP_CONTENT_TYPE'] : '';
if ($contentType === '' && !empty($_SERVER['CONTENT_TYPE'])) {
    $contentType = $_SERVER['CONTENT_TYPE'];
}

$body = null;
$postFields = null; // array para multipart; cURL lo enviará con boundary automático

if (in_array($method, ['POST', 'PUT', 'PATCH', 'DELETE'], true)) {
    $isMultipart = (stripos($contentType, 'multipart/form-data') !== false);

    if ($isMultipart) {
        // Reconstruir cuerpo desde $_FILES y $_POST; php://input no está disponible en multipart.
        $postFields = [];
        foreach ($_FILES as $fieldName => $file) {
            if (is_array($file['error'])) {
                foreach ($file['error'] as $idx => $err) {
                    if ($err === UPLOAD_ERR_OK && !empty($file['tmp_name'][$idx])) {
                        $mime = isset($file['type'][$idx]) ? $file['type'][$idx] : 'application/octet-stream';
                        $postFields[$fieldName . '[' . $idx . ']'] = new CURLFile($file['tmp_name'][$idx], $mime, $file['name'][$idx]);
                    }
                }
            } else {
                if ($file['error'] === UPLOAD_ERR_OK && !empty($file['tmp_name'])) {
                    $mime = isset($file['type']) ? $file['type'] : 'application/octet-stream';
                    $postFields[$fieldName] = new CURLFile($file['tmp_name'], $mime, $file['name']);
                }
            }
        }
        foreach ($_POST as $key => $value) {
            $postFields[$key] = $value;
        }
    } else {
        $body = file_get_contents('php://input');
        if ($body === '') {
            $body = null;
        }
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
if ($postFields === null) {
    // No multipart: enviar Content-Type (cURL lo fija automáticamente cuando postFields es array).
    if (isset($_SERVER['HTTP_CONTENT_TYPE']) && $_SERVER['HTTP_CONTENT_TYPE'] !== '') {
        $headers[] = 'Content-Type: ' . $_SERVER['HTTP_CONTENT_TYPE'];
    } elseif (!empty($_SERVER['CONTENT_TYPE'])) {
        $headers[] = 'Content-Type: ' . $_SERVER['CONTENT_TYPE'];
    } elseif ($body !== null) {
        $headers[] = 'Content-Type: application/json';
    }
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
if ($postFields !== null) {
    curl_setopt($ch, CURLOPT_POSTFIELDS, $postFields);
} elseif ($body !== null) {
    curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
}

$response = curl_exec($ch);
$errno = curl_errno($ch);
$httpCode = (int) curl_getinfo($ch, CURLINFO_HTTP_CODE);
$headerSize = (int) curl_getinfo($ch, CURLINFO_HEADER_SIZE);
curl_close($ch);

if ($errno) {
    http_response_code(502);
    header('Content-Type: application/json');
    echo json_encode([
        'message' => 'No se pudo conectar al servicio de formularios.',
        'error'   => curl_strerror($errno),
    ]);
    exit;
}

// Separar cabeceras y cuerpo usando HEADER_SIZE para evitar casos
// donde hay líneas en blanco iniciales u otros bloques de cabeceras.
$responseHeaders = substr($response, 0, $headerSize);
$responseBody = substr($response, $headerSize);

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
