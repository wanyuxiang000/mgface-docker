package nsenter

/*
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

__attribute__((constructor)) void enter_namespace(void) {
	printf("这个包一旦被引用,它就会在所有的go运行的环境启动之前执行,这样就避免了Go多线程导致的无法进入mnt Namespace的问题.这段程序执行完毕后,Go程序才会执行.\n");
	char *mydocker_pid;
	mydocker_pid = getenv("mydocker_pid");
	 printf("Cgo进行exec实现PID:%s\n", mydocker_pid);
	if (mydocker_pid) {
		printf("got mydocker_pid=%s\n", mydocker_pid);
	} else {
		printf("missing mydocker_pid env skip nsenter\n");
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
		if (setns(fd, 0) == -1) {
			fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
	int res = system(mydocker_cmd);
	exit(0);
	return;
}
int test() {
    return 2016;
}
*/
import "C"
