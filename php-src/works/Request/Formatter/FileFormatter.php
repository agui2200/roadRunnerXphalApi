<?php
/**
 * Created by PhpStorm.
 * User: niebangheng
 * Date: 2019/1/27
 * Time: 19:13
 */

namespace Works\Request\Formatter;


use PhalApi\Exception\BadRequestException;

class FileFormatter extends \PhalApi\Request\Formatter\FileFormatter
{
    public function parse($value, $rule) {
        $default = isset($rule['default']) ? $rule['default'] : NULL;
        $index = $rule['name'];
        $fileList = array();

        // 未上传 && (有默认值 || 非必须)
        if (!isset($_FILES[$index]) && ($default !== NULL || empty($rule['require']))) {
            return $default;
        }
        if (!isset($_FILES[$index])) {
            $message = isset($rule['message'])
                ? \PhalApi\T($rule['message'])
                : \PhalApi\T('miss upload file: {file}', array('file' => $index));
            throw new BadRequestException($message);
        }

        if (is_array($_FILES[$index])) {
            foreach ($_FILES[$index] as $tfile) {
                $file = array(
                    'name' => $tfile->getClientFilename(),
                    'type' => $tfile->getClientMediaType(),
                    'file' => $tfile,
                    'error' => $tfile->getError(),
                    'size' => $tfile->getSize(),
                );

                $fileList[] = $this->parseOne($file, $rule);
            }
        } else {
            $file = array(
                'name' => $_FILES[$index]->getClientFilename(),
                'type' => $_FILES[$index]->getClientMediaType(),
                'file' => $_FILES[$index],
                'error' => $_FILES[$index]->getError(),
                'size' => $_FILES[$index]->getSize(),
            );
            return $file;
        }
        // 返回文件信息二维数组
        return $fileList;
    }

}