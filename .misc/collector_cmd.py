#! /usr/bin/python

import time

def main():
	print "gauge|rpm|{key}|{t:.0f}|{v}".format(key="metrics.cmd.python.1", t=time.time(), v=123)
	print "counter|khz|{key}|{t:.0f}|{v}".format(key="metrics.cmd.python.1", t=time.time(), v=123.1)

if __name__ == '__main__':
	main()
