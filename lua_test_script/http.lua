require "xjson"

response = {charset="utf-8"}

function main(args)
    -- print("http")
    local httpClient = http.new()
    -- print("new http client")
    body, status, err = httpClient:get("http://www.baidu.com", "utf-8")
    -- print("get http request")
    -- if err ~= nil then
    --     print("err")
    --     print(err)
    --     return err
    -- end

    -- if status ~= "200" then
    --     print("status")
    --     print(status)
    --     return status
    -- end

    -- print(body)
    print(status)
    print(err)

    return { status = status, err = err, data = "data1" }
end