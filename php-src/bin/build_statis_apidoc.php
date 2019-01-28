<?php
/**
 * Created by PhpStorm.
 * User: niebangheng
 * Date: 2019/1/26
 * Time: 15:52
 */

// 命令示例：php bin\build_statis_apidoc.php 127.0.0.1 expand
// 主题风格，fold = 折叠，expand = 展开


global $argv;

require_once dirname(__FILE__) . '/../public/init.php';

$projectName = 'PhalApi开源接口框架';

if (!isset($argv[1])) {
    echo "
Usage:

生成展开版：  php {$argv[0]} [接口域名] expand
生成折叠版：  php {$argv[0]} [接口域名] fold
    " . PHP_EOL;
    return;
}


$_SERVER['SCRIPT_NAME'] = '/index.php';
$_SERVER['HTTP_HOST'] = $argv[1];

$theme = isset($argv[2]) ? $argv[2] : 'fold';

$apiList = new \PhalApi\Helper\ApiStaticCreate($projectName, $theme);
$apiList->render();