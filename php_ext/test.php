<?php
dl('mcrypt.so');
dl('jindan.so');


function test($data) {
    $etype = MCRYPT_RIJNDAEL_256;
    $iv    = mcrypt_create_iv(mcrypt_get_iv_size($etype, MCRYPT_MODE_ECB), MCRYPT_RAND);

    //和mcrypt_encrypt($etype, $secret_key, $string, MCRYPT_MODE_CBC, $iv)一样,只是把$secret_key写死在扩展里
    $en = walu_hello($etype, $data, MCRYPT_MODE_CBC, $iv);
    echo "\n加密:" . BIN2HEX($en);

    //$de = decypt_bank($etype, $en, MCRYPT_MODE_CBC, $iv);
    //echo "\n解密:$de";
}

test('aaa');
