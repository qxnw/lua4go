require "xtable"

response = {charset="utf-8"}

function main(args)
    print("response")
    print("session:" )
    print(__session__)
    print("loggerName:" )
    print(__logger_name__)
    print("http_context")
    print(__http_context__)

    -- input = args.input
    print("args的类型")
    print(type(args))
    -- print("内容")
    -- print(xtable.tojson(args))
    return args
end