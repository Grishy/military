$(function () {
    $('#tree').jstree({
        'core': {
            'data': {
                'url': '/tree/get_node',
                'data': function (node) {
                    return { 'id': node.id };
                }
            },
            'force_text': true,
            'check_callback': true,
            'themes': {
                'responsive': true
            }
        },
        'plugins': ['state', 'dnd', 'contextmenu', 'wholerow'],
        "contextmenu":{
            'items' : function(node) {
                var items = $.jstree.defaults.contextmenu.items();
                items.ccp = false;

                return items;
            }
        }
    })
        .on('delete_node.jstree', function (e, data) {
            $.get('/tree/delete_node', { 'id': data.node.id })
                .fail(function () {
                    data.instance.refresh();
                });
        })
        .on('create_node.jstree', function (e, data) {
            $.get('/tree/create_node', { 'id': data.node.parent, 'position': data.position, 'text': data.node.text })
                .done(function (d) {
                    data.instance.set_id(data.node, d.id);
                })
                .fail(function () {
                    data.instance.refresh();
                });
        })
        .on('rename_node.jstree', function (e, data) {
            $.get('/tree/rename_node', { 'id': data.node.id, 'text': data.text })
                .fail(function () {
                    data.instance.refresh();
                });
        })
        .on('move_node.jstree', function (e, data) {
            $.get('/tree/move_node', { 'id': data.node.id, 'parent': data.parent, 'position': data.position })
                .fail(function () {
                    data.instance.refresh();
                });
        })
        .on('changed.jstree', function (e, data) {
            if (data && data.selected && data.selected.length) {
                $.get('/tree/get_content?id=' + data.selected.join(':'), function (d) {
                    $('.m-main').text(d.content).show();
                });
            }
            else {
                $('#data .content').hide();
                $('.m-main').text('Select a file from the tree.').show();
            }
        });
});