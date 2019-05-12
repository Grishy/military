function getQueryVariable(variable)
{
       var query = window.location.search.substring(1);
       var vars = query.split("&");
       for (var i=0;i<vars.length;i++) {
               var pair = vars[i].split("=");
               if(pair[0] == variable){return pair[1];}
       }
       return(false);
}

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
        'plugins': [ 'dnd', 'contextmenu', 'wholerow', "types"],
        "types": {
            "default": {
                "icon": "glyphicon glyphicon-file"
            },
        },
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

                    if (d.content) {
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


    setTimeout(function () {
        getQueryVariable()
        $('#tree')
            .jstree('open_node', window.location.search.substring(1));
    }, 500)
    // Initialize our editor
    var editor = ContentTools.EditorApp.get();
    editor.init('*[data-editable]', 'data-name');

    editor.isReady(console.log)

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
        payload.append("text", regions["formatted"]);

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

    ContentTools.IMAGE_UPLOADER = function (dialog) {
        var imagePath, imageSize, rotate, uploadingTimeout;
        // return imagePath = "/images/pages/demo/landscape-in-eire.jpg",
        // imageSize = [780, 366],
        uploadingTimeout = null;

        function rotate() {
            var clearBusy;
            return dialog.busy(!0),
                clearBusy = function (_this) {
                    return function () {
                        return dialog.busy(!1)
                    }
                }(this),
                setTimeout(clearBusy, 1500)
        };

        dialog.addEventListener("imageuploader.cancelupload", function () {
            // Cancel the current upload

            // Stop the upload
            if (xhr) {
                xhr.upload.removeEventListener('progress', xhrProgress);
                xhr.removeEventListener('readystatechange', xhrComplete);
                xhr.abort();
            }

            // Set the dialog to empty
            dialog.state('empty');
        });

        dialog.addEventListener("imageuploader.clear", function () {
            return dialog.clear()
        });

        dialog.addEventListener('imageuploader.fileready', function (ev) {
            // Upload a file to the server
            var formData;
            var file = ev.detail().file;

            // Define functions to handle upload progress and completion
            xhrProgress = function (ev) {
                // Set the progress for the upload
                dialog.progress((ev.loaded / ev.total) * 100);
            }

            xhrComplete = function (ev) {
                var response;

                // Check the request is complete
                if (ev.target.readyState != 4) {
                    return;
                }

                // Clear the request
                xhr = null
                xhrProgress = null
                xhrComplete = null

                // Handle the result of the upload
                if (parseInt(ev.target.status) == 200) {
                    // Unpack the response (from JSON)
                    response = JSON.parse(ev.target.responseText);

                    // Store the image details
                    image = {
                        size: response.size,
                        url: response.url
                    };

                    // Populate the dialog
                    dialog.populate(image.url, image.size);

                } else {
                    // The request failed, notify the user
                    new ContentTools.FlashUI('no');
                }
            }

            // Set the dialog state to uploading and reset the progress bar to 0
            dialog.state('uploading');
            dialog.progress(0);

            // Build the form data to post to the server
            formData = new FormData();
            formData.append('image', file);

            // Make the request
            xhr = new XMLHttpRequest();
            xhr.upload.addEventListener('progress', xhrProgress);
            xhr.addEventListener('readystatechange', xhrComplete);
            xhr.open('POST', '/upload-image', true);
            xhr.send(formData);
        });

        dialog.addEventListener("imageuploader.rotateccw", function () {
            return rotate()
        });

        dialog.addEventListener("imageuploader.rotatecw", function () {
            return rotate()
        });

        dialog.addEventListener("imageuploader.save", function () {
            var crop, cropRegion, formData;

            // Define a function to handle the request completion
            xhrComplete = function (ev) {
                // Check the request is complete
                if (ev.target.readyState !== 4) {
                    return;
                }

                // Clear the request
                xhr = null
                xhrComplete = null

                // Free the dialog from its busy state
                dialog.busy(false);

                // Handle the result of the rotation
                if (parseInt(ev.target.status) === 200) {
                    // Unpack the response (from JSON)
                    var response = JSON.parse(ev.target.responseText);

                    // Trigger the save event against the dialog with details of the
                    // image to be inserted.
                    dialog.save(
                        response.url,
                        response.size,
                        {
                            'alt': response.alt,
                            'data-ce-max-width': response.size[0]
                        });

                } else {
                    // The request failed, notify the user
                    new ContentTools.FlashUI('no');
                }
            }

            // Set the dialog to busy while the rotate is performed
            dialog.busy(true);

            // Build the form data to post to the server
            formData = new FormData();
            formData.append('url', image.url);

            // Set the width of the image when it's inserted, this is a default
            // the user will be able to resize the image afterwards.
            formData.append('width', 600);

            // Check if a crop region has been defined by the user
            if (dialog.cropRegion()) {
                formData.append('crop', dialog.cropRegion());
            }

            // Make the request
            xhr = new XMLHttpRequest();
            xhr.addEventListener('readystatechange', xhrComplete);
            xhr.open('POST', '/insert-image', true);
            xhr.send(formData);
        })
    }

});