

box.cfg {
    listen = 3301,
    log_level = 5,
    memtx_memory = 128 * 1024 * 1024, 
    memtx_min_tuple_size = 16,
    memtx_max_tuple_size = 10 * 1024 * 1024, 
    wal_mode = "write",
    wal_dir = "/var/lib/tarantool",
    memtx_dir = "/var/lib/tarantool",
    log = "/var/log/tarantool.log",  
    readahead = 10 * 1024 * 1024 
}

local function create_polls_space()
    local s = box.schema.space.create('polls', {
        if_not_exists = true,
        format = {
            {name = 'id', type = 'string'},
            {name = 'data', type = 'map'}
        }
    })
    
    s:create_index('primary', {
        type = 'HASH',
        parts = {1, 'string'},
        if_not_exists = true
    })
    
    print('Polls space created successfully')
end


box.once('access:voting_app', function()
    box.schema.user.create('voting_app', {
        if_not_exists = true,
        password = 'voting_app_password'
    })
    
    box.schema.user.grant('voting_app', 'read,write,execute', 'universe', nil, {if_not_exists = true})
    print('User voting_app created and granted permissions')
end)


create_polls_space()

print('Tarantool initialization completed')


require('console').start()
os.execute("chmod 777 /var/log/tarantool.log")