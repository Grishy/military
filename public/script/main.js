$(function () {
    let currID;

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
        "contextmenu": {
            'items': function (node) {
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
                    currID = data.node.id
                    
                    $('.m-main-title').val(data.node.text);
                    if (!d.content) {
                        $('.m-main-text').html(d.content);
                    } else {
                        $('.m-main-text').html("");
                    }
                });
            }
            else {
                $('.m-main-text').text('Выбериту нужный раздел.').show();
            }
        });

    // Initialize our editor
    var editor = ContentTools.EditorApp.get();
    editor.init('*[data-editable]', 'data-name');

    editor.addEventListener('saved', function (ev) {
        var name, onStateChange, passive, payload, regions, xhr;

        // Check if this was a passive save
        passive = ev.detail().passive;

        // Check to see if there are any changes to save
        regions = ev.detail().regions;

        // Set the editors state to busy while we save our changes
        this.busy(true);

        // Collect the contents of each region into a FormData instance
        payload = new FormData();
        payload.append('id', currID);
        payload.append('title', $(".m-main-title").val());
        payload.append("text", regions["main-content"]);

        // Send the update content to the server to be saved
        onStateChange = function (ev) {
            // Check if the request is finished
            if (ev.target.readyState == 4) {
                editor.busy(false);
                if (ev.target.status == '200') {
                    // Save was successful, notify the user with a flash
                    if (!passive) {
                        new ContentTools.FlashUI('ok');
                    }
                } else {
                    // Save failed, notify the user with a flash
                    new ContentTools.FlashUI('no');
                }
            }
        };

        xhr = new XMLHttpRequest();
        xhr.addEventListener('readystatechange', onStateChange);
        xhr.open('POST', '/save-page');
        xhr.send(payload);
    });

});