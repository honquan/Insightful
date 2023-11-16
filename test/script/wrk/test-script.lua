request = function()
  wrk.headers["Connection"] = "Keep-Alive"
  id = math.random(1,100)
  wrk.body = '{"page_id":'..id..', "positions": [{"x": '..math.random(1,1080)..', "y": '..math.random(1,1920)..'}, {"x": '..math.random(1,1080)..', "y": '..math.random(1,1920)..'}, {"x": '..math.random(1,1080)..', "y": '..math.random(1,1920)..'}]}'
  wrk.headers["Content-Type"] = "application/json"
  return wrk.format("POST")
end