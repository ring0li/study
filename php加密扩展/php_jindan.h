#ifndef JINDAN_H
#define JINDAN_H

#define ENCYPT_SECRET "1234567891234567"//长度必须16位 24位 32位
#ifdef HAVE_CONFIG_H
#include "config.h"
#endif

#include "php.h"

#define phpext_jindan_ptr &jindan_module_entry
extern zend_module_entry jindan_module_entry;
#endif
