# Add task
```
curl 
-i -XPOST localhost:18000/tasks -d @./handler/tes
tdata/add_task/ok_req.json.golden
```

# Get task
```
curl -i -XGET localhost:18000/tasks
```