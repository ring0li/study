#ifndef JINDAN_H
#define JINDAN_H

#define ENCYPT_SECRET "jindanlicai.com"

#ifdef HAVE_CONFIG_H
#include "config.h"
#endif

#include "php.h"
#define phpext_jindan_ptr &jindan_module_entry
extern zend_module_entry jindan_module_entry;

#endif
