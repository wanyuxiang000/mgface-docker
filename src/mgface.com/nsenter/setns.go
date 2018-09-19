package nsenter

/*
#define _GNU_SOURCE
#include <unistd.h>
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
//这里的__attribute__((constructor))指的是， 一旦这个包被引用，那么这个函数就会被自动执行
//类似于构造函数，会在程序一启动的时候运行
__attribute__((constructor)) void enter_namespace(void) {
	char *mydocker_pid;
	//从环境变量中获取需要进入的PID
	mydocker_pid = getenv("mydocker_pid");
	if (mydocker_pid) {
		fprintf(stdout, "got mydocker_pid=%s\n", mydocker_pid);
	} else {
	    //这里如果没有指定PID就不需要向下执行,直接退出
		fprintf(stdout, "missing mydocker_pid env skip nsenter");
		return;
	}
	char *mydocker_cmd;
	mydocker_cmd = getenv("mydocker_cmd");
	if (mydocker_cmd) {
		fprintf(stdout, "got mydocker_cmd=%s\n", mydocker_cmd);
	} else {
		fprintf(stdout, "missing mydocker_cmd env skip nsenter");
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };

	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", mydocker_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);
		//这里才真正调用 setns 系统调用进入对应的 Namespace
		if (setns(fd, 0) == -1) {
			fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
	//在进入的 Namespace 中执行指定的命令
	int res = system(mydocker_cmd);
	exit(0);
	return;
}
*/
import "C"

//这个包一旦被引用，它就会在所有的go运行的环境启动之前执行,这样就避免了Go多线程导致的无法进入mnt Namespace的问题。这段程序执行完毕后 ，Go程序才会执行。
