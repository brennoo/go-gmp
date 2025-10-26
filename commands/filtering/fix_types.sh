#!/bin/bash

# Fix integer type issues in options.go
sed -i '' 's/fb\.AddCondition("\([^"]*\)", Op\([^,]*\), \([^)]*\))/fb.AddCondition("\1", Op\2, strconv.Itoa(\3))/g' options.go

# Fix float type issues in options.go  
sed -i '' 's/fb\.AddCondition("\([^"]*\)", Op\([^,]*\), \([^)]*\))/fb.AddCondition("\1", Op\2, fmt.Sprintf("%.1f", \3))/g' options.go

# Fix boolean type issues in options.go
sed -i '' 's/fb\.AddCondition("\([^"]*\)", Op\([^,]*\), \([^)]*\))/fb.AddCondition("\1", Op\2, strconv.FormatBool(\3))/g' options.go
