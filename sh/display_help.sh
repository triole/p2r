#!/bin/bash

r -h | sed -E 's|(config=").*?(")|\1$(pwd)/p2r.yaml\2|g'
