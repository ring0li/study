#include "php_jindan.h"

//功能实现,加密
ZEND_FUNCTION(encypt_bank) {
    char *cipher, *data, *mode, *IV;
    int cipher_len, data_len, mode_len, IV_len;

    zval retval, funcname;




    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "s|s|s|s", &cipher, &cipher_len, &data, &data_len, &mode, &mode_len, &IV, &IV_len) == FAILURE){
        WRONG_PARAM_COUNT;
        return;
    }

    zend_uint param_count = ZEND_NUM_ARGS();
    zval *args[param_count];

    MAKE_STD_ZVAL(args[0]);
    MAKE_STD_ZVAL(args[1]);
    MAKE_STD_ZVAL(args[2]);
    MAKE_STD_ZVAL(args[3]);
    MAKE_STD_ZVAL(args[4]);

//    ZVAL_STRINGL(&funcname, "mcrypt_encrypt", sizeof("mcrypt_encrypt") - 1, 1);
    ZVAL_STRING(&funcname, "mcrypt_encrypt", 0);
//    ZVAL_STRINGL(&funcname, "var_dump", sizeof("var_dump") - 1, 1);
    ZVAL_STRINGL(args[0], cipher, cipher_len, 1);
    ZVAL_STRINGL(args[1], "1", sizeof("1") - 1, 1);
    ZVAL_STRINGL(args[2], data, data_len, 1);
    ZVAL_STRINGL(args[3], mode, mode_len, 1);
    ZVAL_STRINGL(args[4], IV, IV_len, 1);

    PHPWRITE(cipher, cipher_len);
    printf("cipher_len:%d\n", cipher_len);
    php_printf("\n");
    PHPWRITE(data, data_len);

//    TSRMLS_FETCH();
    if (call_user_function(EG(function_table), NULL, &funcname, &retval, 5, args TSRMLS_CC) == FAILURE) {
        RETURN_BOOL(0);//调用失败
    }
    else
    {
        RETVAL_STRING(Z_STRVAL(retval), 1);
    }

    printf("cipher_len:%d\n", cipher_len);
    zval_dtor(args[0]);
    zval_dtor(args[1]);
    zval_dtor(args[2]);
    zval_dtor(args[3]);
    zval_dtor(args[4]);
//    zval_ptr_dtor(&funcname);
    zval_dtor(&retval);
    printf("cipher_len:%d\n", cipher_len);
}

//功能list
const zend_function_entry jindan_functions[] = {
        ZEND_FE(encypt_bank, NULL)
        {NULL, NULL, NULL}
};

//引入到zend中
zend_module_entry jindan_module_entry = {
        STANDARD_MODULE_HEADER,
        "jindan",
        jindan_functions,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        "0.1",
        STANDARD_MODULE_PROPERTIES
};

//动态加载,类似于dl()函数
#ifdef COMPILE_DL_JINDAN
ZEND_GET_MODULE(jindan)
#endif
