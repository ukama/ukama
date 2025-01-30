#/bin/bash
SCRIPT_DIR=$1
cat /proc/mounts | cut -d' ' -f2 | grep "^$SCRIPT_DIR." | sort -r | while read path; do
			echo "Unmounting $path" >&2
			$_sudo umount -fn "$path" || exit 1
		done
