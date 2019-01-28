<?php


/**
 * @var Goridge\RelayInterface $relay
 */

use Spiral\Goridge;
use Spiral\RoadRunner;
use Works\Response\JsonResponse;

ini_set('display_errors', 'stderr');

require 'init.php';

$worker = new RoadRunner\Worker(new Goridge\StreamRelay(STDIN, STDOUT));
$psr7 = new RoadRunner\PSR7Client($worker);

$di = \PhalApi\DI();

while ($req = $psr7->acceptRequest()) {
    global $response;
    $response = new \Zend\Diactoros\Response();
    parse_str($req->getUri()->getQuery(), $_GET);
    $_POST = $req->getParsedBody();
    $_FILES = $req->getUploadedFiles();
    $urlPath = explode('/', $req->getUri()->getPath());
    $firstPath = isset($urlPath[1]) ? $urlPath[1] : '';
    $di->response = new JsonResponse();
    $di->request = new Works\Request\PsrRequest(array_merge((array)$_GET, (array)$_POST));
    $di->_formatterFile = new \Works\Request\Formatter\FileFormatter();
    ob_start();
    try {
        switch ($firstPath) {
            case 'docs':
                $_SERVER['REQUEST_URI'] = '/docs';
                $_SERVER['SCRIPT_NAME'] = '/index.php';
                if (isset($_GET['detail']) && $_GET['detail']) {
                    $api = new \PhalApi\Helper\ApiDesc('接口文档');
                    $api->render();
                } else {
                    $api = new \PhalApi\Helper\ApiList('接口文档');
                    $api->render();
                }
                break;
            case 'api':
                $pai = new \PhalApi\PhalApi();
                $pai->response()->output();
                break;
        }
        $buffer = ob_get_contents();
        $response->getBody()->write($buffer);
        $psr7->respond($response);
    } catch (\Throwable $e) {
        $psr7->getWorker()->error((string)$e);
    }
    ob_clean();
}