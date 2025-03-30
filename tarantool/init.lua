box.cfg{
    listen = 3301
}

box.schema.space.create('votes', {
    if_not_exists = true,
    format = {
        {name = 'id', type = 'unsigned'}, 
        {name = 'channel_id', type = 'string'}, 
        {name = 'options', type = 'array'},
        {name = 'votes', type = 'map'},
        {name = 'active', type = 'boolean'}
    }
})

box.space.votes:create_index('primary', {
    type = 'tree',
    parts = {'id'},
    if_not_exists = true
})

box.schema.user.create('guest', {if_not_exists = true})
box.schema.user.grant('guest', 'read,write,execute', 'universe', nil, {if_not_exists = true})