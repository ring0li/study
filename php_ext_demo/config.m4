PHP_ARG_ENABLE(jindan,
[Whether to enable the "jindan" extension],
[    enable-jindan Enable "jindan" extension support])

if test $PHP_JINDAN != "no"; then
  PHP_SUBST(JINDAN_SHARED_LIBADD)
  PHP_NEW_EXTENSION(jindan, jindan.c, $ext_shared)
fi

if test -z "$PHP_DEBUG"; then
        AC_ARG_ENABLE(debug,
                [--enable-debug  compile with debugging system],
                [PHP_DEBUG=$enableval], [PHP_DEBUG=no]
        )
fi
