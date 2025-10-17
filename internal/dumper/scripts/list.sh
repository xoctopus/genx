#!/bin/sh

go list std | grep -v -E 'vendor($|/)|internal' > std.list