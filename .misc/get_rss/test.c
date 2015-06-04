/*
 * test.c
 * Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
 *
 * Distributed under terms of the MIT license.
 */

#include <stdio.h>
#include "get_rss.c"

int main(){
  printf("PeakRSS: %dK, CurrentRSS: %dK\n", getPeakRSS()/1024, getCurrentRSS()/1024);
}

