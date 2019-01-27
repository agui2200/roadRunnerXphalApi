<?php
/**
 * Created by PhpStorm.
 * User: niebangheng
 * Date: 2019/1/27
 * Time: 18:33
 */

namespace Works\Response;

use PhalApi;
use Zend;

abstract class Response extends PhalApi\Response
{
    protected function handleHeaders($headers) {
        global $response;
        foreach ($headers as $key => $content) {
            $response = $response->withAddedHeader($key, $content);
        }
    }

}