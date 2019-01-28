<?php
$fd = fopen('log','w+');
fwrite($fd,time()."\r\n");
fclose($fd);