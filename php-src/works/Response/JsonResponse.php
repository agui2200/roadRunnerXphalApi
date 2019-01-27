<?php
/**
 * Created by PhpStorm.
 * User: niebangheng
 * Date: 2019/1/27
 * Time: 18:33
 */

namespace Works\Response;

class JsonResponse extends Response
{
    /**
     * @var int JSON常量组合的二进制掩码
     * @see http://php.net/manual/en/json.constants.php
     */
    protected $options;

    public function __construct($options = 0) {
        $this->options = $options;

        $this->addHeaders('Content-Type', 'application/json;charset=utf-8');
    }

    protected function formatResult($result) {
        return json_encode($result, $this->options);
    }
}