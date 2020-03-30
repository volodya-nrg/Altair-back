$(function () {
    $(".accordion").accordion({collapsible: true, heightStyle: "content", active: 100});
    $(".datepicker").datepicker({dateFormat: "yy-mm-dd"});

    insertCatsTreeAsTagSelect('.target_for_cats_tree');
    insertKindsProperties('.target_for_kind_properties');
    insertProperties('.target_for_properties');
    insertValuesForProperties('.target_for_values_properties');

    $(document).on('click', '.wrapper_for_photo .icon', function (e) {
        var $wrapperForPhoto = $(this).closest('.wrapper_for_photo');
        $wrapperForPhoto.remove();
    });
    $(document).on('click', '.target_for_properties .dynamic__add', function (e) {
        addProperty($(this));
    });
    $(document).on('click', '.target_for_properties .dynamic__del', removeProperty);
    $(document).on('click', '.target_for_values_properties .dynamic__add', function (e) {
        addValueForProperty($(this));
    });
    $(document).on('click', '.target_for_values_properties .dynamic__del', delValueForProperty);
    $(document).on('change', '.form[action="/api/v1/ads"][method="post"] .target_for_cats_tree > select, .form-put-put-ads .target_for_cats_tree > select', function (e) {
        changeSelectOnCatsTree($(e.target));
    });

    $(".form").on("submit", function (e) {
        e.preventDefault();

        var $form = $(this);
        var $submit = $form.find("[type='submit']");
        var $jsonResult = $('.json-result');
        var action = getFineAction($form);
        var isFormPutGetUsers = $form.hasClass('form-put-get-users');
        var isFormPutGetCats = $form.hasClass('form-put-get-cats');
        var isFormPutGetAds = $form.hasClass('form-put-get-ads');
        var isFormPutGetProperties = $form.hasClass('form-put-get-properties');
        var isFormPutGetKindProperties = $form.hasClass('form-put-get-kind_properties');
        var data = new FormData(this);
        var objSettings = {
            processData: false,  // Important!
            contentType: false,
            cache: false,
            method: $form.attr('method'),
            url: action,
            data: data,
            dataType: 'json',
            beforeSend: function (xhr) {
                $submit.attr("disabled", true);
                $jsonResult.text('');
            }
        };

        if ($form.attr('enctype')) {
            objSettings.enctype = $form.attr('enctype');
        }

        $.ajax(objSettings).done(function (response) {
            var text = JSON.stringify(response, null, '\t');

            $jsonResult.text(text);

            if (isFormPutGetUsers) {
                $tmpForm = $('.form-put-put-users');
                $tmpForm.addClass("hidden").get(0).reset();
                $tmpForm.find('.form__files').empty();
                formPutPutUsers(response, $tmpForm);

            } else if (isFormPutGetCats) {
                $tmpForm = $('.form-put-put-cats');
                $tmpForm.addClass("hidden").get(0).reset();
                clickDynamicDel($tmpForm);
                formPutPutCats(response, $tmpForm);

            } else if (isFormPutGetAds) {
                $tmpForm = $('.form-put-put-ads');
                $tmpForm.addClass("hidden").get(0).reset();
                $tmpForm.find('.form__files').empty();
                $tmpForm.find('#catsTree').empty();
                formPutPutAds(response, $tmpForm);

            } else if (isFormPutGetProperties) {
                $tmpForm = $('.form-put-put-properties');
                $tmpForm.addClass("hidden").get(0).reset();
                clickDynamicDel($tmpForm);
                formPutPutProperties(response, $tmpForm);

            } else if (isFormPutGetKindProperties) {
                $tmpForm = $('.form-put-put-kind_properties');
                $tmpForm.addClass("hidden").get(0).reset();
                formPutPutKindProperties(response, $tmpForm);
            }

            if ($form.hasClass('sx-reload')) {
                window.location.reload();
            }

        }).fail(function (data) {
            alert("Status: " + data.status + "; " + data.responseText);

        }).always(function () {
            $submit.attr("disabled", false);
        });
    });
});

function appendPhotos($form, name, aImages) {
    $files = $form.find('.form__files');

    if (!$files.length) {
        return
    }

    for (var i = 0; i < aImages.length; i++) {
        var img = aImages[i];
        var tpl = [
            '<div class="wrapper_for_photo">',
            '<span class="icon">x</span>',
            '<img height="30" src="/images/' + img.filepath + '"/>',
            '<input type="hidden" name="' + name + '" value="' + img.filepath + '"/>',
            '</div>',
        ].join('');
        $files.append(tpl);
    }
}

function insertCatsTreeAsTagSelect(selectorTarget) {
    var $targetts = $(selectorTarget);
    var $select = $('<select></select>');

    walkForCatsTree(ALTAIR.catsTree.childes, $select);

    $targetts.each(function () {
        var $self = $(this);
        var name = $self.data("name");
        var $selectCopy = $select.clone();

        $selectCopy.attr("name", name);
        $selectCopy.prepend('<option value=""></option>');
        $selectCopy.val(0);

        $self.append($selectCopy);
    });

    function walkForCatsTree(branch, $reciever, prefixSrc) {
        var prefix = prefixSrc || "";

        for (var key in branch) {
            var el = branch[key];

            $reciever.append('<option value="' + el.catId + '">' + prefix + el.name + '</option>');

            if (el.childes && el.childes.length) {
                walkForCatsTree(el.childes, $reciever, prefix + "|----");
            }
        }
    }
}

function insertKindsProperties(selectorPlace) {
    var $targetPlace = $(selectorPlace);
    var $select = $('<select></select>');

    for (var key in ALTAIR.kindProperties) {
        var el = ALTAIR.kindProperties[key];
        $select.append('<option value="' + el.kindPropertyId + '">' + el.name + '</option>');
    }

    $targetPlace.each(function () {
        var $self = $(this);
        var name = $self.data("name");
        var $selectCopy = $select.clone();

        $selectCopy.attr("name", name);
        $self.append($selectCopy);
    });
}

function insertProperties(selectorPlace) {
    var tpl = [
        '<div class="dynamic">',
        '   <div class="dynamic__controls">',
        '       <div><span class="icon dynamic__add">+</span></div>',
        '   </div>',
        '   <div class="dynamic__items"></div>',
        '</div>',
    ].join('');
    var $dynamic = $(tpl);
    var $select = $('<select></select>');

    $select.append('<option value="0"></option>');
    for (var key in ALTAIR.properties) {
        var el = ALTAIR.properties[key];
        $select.append('<option value="' + el.propertyId + '" data-title="' + el.title + '">' + el.title + '</option>');
    }

    $(selectorPlace).each(function () {
        var $self = $(this);
        var $dynamicCopy = $dynamic.clone();

        $dynamicCopy.find('.dynamic__controls').prepend($select.clone());
        $self.append($dynamicCopy);
    });
}

function insertValuesForProperties(selector) {
    var $places = $(selector);
    var $dynamic = $([
        '<div class="dynamic">',
        '   <div class="dynamic__controls">',
        '       <div></div>',
        '       <div><span class="icon dynamic__add">+</span></div>',
        '   </div>',
        '   <div class="dynamic__items"></div>',
        '</div>',
    ].join(''));

    $places.each(function () {
        var $self = $(this);
        $self.append($dynamic.clone());
    });
}

function clickDynamicDel($ctx) {
    $ctx.find(".dynamic__del").each(function () {
        $(this).trigger("click");
    });
}

function getFineAction($form) {
    var url = $form.attr('action');
    var method = $form.attr("method") || "get";
    var aFormData = $form.serializeArray();
    var serialize = $form.serialize();
    var indexOfQues = url.indexOf("?");
    var path = (indexOfQues !== -1 ? url.substr(0, indexOfQues) : url);
    var sQuery = (indexOfQues !== -1 ? url.substr(indexOfQues + 1) : "");
    var aParts = path.split("/");
    var result = "";

    for (var i = 0; i < aParts.length; i++) {
        var part = aParts[i];

        if (part === "") {
            continue;
        }
        if (part.charAt(0) === ":") {
            var tmpPart = part.slice(1);

            for (var j = 0; j < aFormData.length; j++) {
                var oItem = aFormData[j];

                if (oItem.name === tmpPart) {
                    part = oItem.value;
                    break;
                }
            }
        }

        result += "/" + part;
    }

    if (serialize && method === "get") {
        if (sQuery) {
            sQuery += "&"
        }
        sQuery += serialize;
    }
    if (sQuery) {
        result += "?" + sQuery;
    }

    return result;
}

function formPutPutUsers(data, $form) {
    $form.removeClass('hidden');
    $form.find('input[name="userId"]').val(data.userId);
    $form.find('input[name="name"]').val(data.name);
    $form.find('input[name="email"]').val(data.email);
    $form.find('input[name="emailIsConfirmed"]').prop("checked", data.emailIsConfirmed);

    if (data.avatar) {
        // storage.Image
        var objs = [
            {
                img_id: 0,
                el_id: 0,
                is_disabled: false,
                opt: "avatar",
                created_at: "",
                filepath: data.avatar,
            }
        ];
        appendPhotos($form, "avatar", objs);
    }
}

function formPutPutCats(data, $form) {
    $form.removeClass('hidden');

    $form.find('input[name="catId"]').val(data.catId);
    $form.find('input[name="name"]').val(data.name);
    $form.find('input[name="slug"]').val(data.slug);
    $form.find('select[name="parentId"]').val(data.parentId);
    $form.find('input[name="pos"]').val(data.pos);
    $form.find('input[name="isDisabled"]').prop("checked", data.isDisabled);

    if (data.properties && data.properties.length) {
        for (var i = 0; i < data.properties.length; i++) {
            addProperty($form.find('.dynamic__add'), data.properties[i]);
        }
    }
}

function formPutPutAds(data, $form) {
    var $selectCat = $form.find('select[name="catId"]');

    $form.removeClass('hidden');
    $form.find('input[name="adId"]').val(data.adId);
    $form.find('input[name="title"]').val(data.title);
    $form.find('input[name="slug"]').val(data.slug);
    $form.find('input[name="userId"]').val(data.userId);
    $form.find('textarea[name="text"]').val(data.text);
    $form.find('input[name="price"]').val(data.price);
    $form.find('input[name="isDisabled"]').prop("checked", data.isDisabled);
    $selectCat.val(data.catId);

    if (data.images.length) {
        appendPhotos($form, "filesAlreadyHas[]", data.images);
    }

    // в пришедшие позже св-ва вставим актуальные
    changeSelectOnCatsTree($selectCat, function () {
        var $box = $form.find('.cat_properties');

        for (var i = 0; i < data.details.length; i++) {
            var detail = data.details[i];
            var kind = detail.kindPropertyName;

            if (kind === "radio") {
                $box.find('input[type="radio"][name="' + detail.propertyName + '"][value="' + detail.value + '"]').prop("checked", true);

            } else {
                $box.find('[name="' + detail.propertyName + '"]').val(detail.value);
            }
        }
    });
}

function formPutPutProperties(data, $form) {
    $form.removeClass('hidden');
    $form.find('input[name="propertyId"]').val(data.propertyId);
    $form.find('input[name="title"]').val(data.title);
    $form.find('input[name="name"]').val(data.name);
    $form.find('input[name="maxInt"]').val(data.maxInt);
    $form.find('select[name="kindPropertyId"]').val(data.kindPropertyId);
    $form.find('input[name="isCanAsFilter"]').prop("checked", data.isCanAsFilter);

    if (data.values && data.values.length) {
        for (var i = 0; i < data.values.length; i++) {
            addValueForProperty($form.find('.dynamic__add'), data.values[i]);
        }
    }
}

function formPutPutKindProperties(data, $form) {
    $form.removeClass('hidden');
    $form.find('input[name="kindPropertyId"]').val(data.kindPropertyId);
    $form.find('input[name="name"]').val(data.name);
}

function addProperty($ctx, properties) {
    var $owner = $ctx;

    if (!$ctx.hasClass('dynamic')) {
        $owner = $ctx.closest('.dynamic');
    }

    var $items = $owner.find('.dynamic__items');
    var $item = $owner.find('.dynamic__item');
    var $select = $owner.find('.dynamic__controls select');
    var propertyIdSrc = parseInt($select.val());
    var index = $item.length + 1;
    var propertyId = properties && properties.propertyId || propertyIdSrc;
    var pos = properties && properties.propertyPos || index;
    var propertyIsRequire = properties && properties.propertyIsRequire || false;

    if (!propertyId) {
        alert('Ошибка: не выбрано значение!');
        return;
    }

    var $option = $select.find('option[value="' + propertyId + '"]');
    var tpl = [
        '<div class="dynamic__item" data-property_id="' + propertyId + '">',
        '   <input type="hidden" name="propertyId[' + index + ']" value="' + propertyId + '"/>',
        '   <div class="dynamic__inputs">',
        '       <input type="text" value="' + $option.data('title') + '" readonly="readonly"/>',
        '       <input class="dynamic__input_mid" type="number" name="pos[' + index + ']" value="' + pos + '"/>',
        '       <label class="dynamic__input_lg">',
        '           <input type="checkbox" name="isRequire[' + index + ']" value="true"' + (propertyIsRequire ? ' checked="checked"' : "") + '/> обяз.',
        '       </label>',
        '       <div class="dynamic__input_short"><span class="icon dynamic__del">-</span></div>',
        '   </div>',
        '</div>',
    ].join('');

    $select.find('option[value="' + propertyId + '"]').attr("disabled", true);
    $items.append(tpl);
    $select.val(0);
}

function removeProperty(e) {
    var $self = $(e.target);
    var $parent = $self.closest('.dynamic__item');
    var $owner = $self.closest('.dynamic');
    var $select = $owner.find('.dynamic__controls select');
    var propertyId = $parent.data('property_id');

    $select.find('option[value="' + propertyId + '"]').attr("disabled", false);
    $parent.remove();
}

function addValueForProperty($ctx, oValue) {
    var $owner = $ctx;

    if (!$ctx.hasClass('dynamic')) {
        $owner = $ctx.closest('.dynamic');
    }

    var $items = $owner.find('.dynamic__items');
    var $item = $owner.find('.dynamic__item');
    var index = $item.length + 1;
    var id = oValue && oValue.valueId || 0;
    var title = oValue && oValue.title || "";
    var pos = oValue && oValue.pos || index;

    var tpl = [
        '<div class="dynamic__item">',
        '   <input type="hidden" name="valueId[' + index + ']" value="' + id + '"/>',
        '   <div class="dynamic__inputs">',
        '       <input type="text" name="valueTitle[' + index + ']" value="' + title + '"/>',
        '       <input class="dynamic__input_mid" type="number" name="valuePos[' + index + ']" value="' + pos + '"/>',
        '       <div class="dynamic__input_short"><span class="icon dynamic__del">-</span></div>',
        '   </div>',
        '</div>',
    ].join('');

    $items.append(tpl);
}

function delValueForProperty(e) {
    var $self = $(e.target);
    var $parent = $self.closest('.dynamic__item');
    $parent.remove();
}

function changeSelectOnCatsTree($select, cb) {
    var catId = $select.val();
    var $wrapper = $('<div class="cat_properties"></div>');

    $.ajax({
        method: 'get',
        url: '/api/v1/cats/' + catId,
        dataType: 'json',
        beforeSend: function (xhr) {
            $select.attr("disabled", true);
            $select.parent().children(".cat_properties").remove();
        }
    }).done(function (response) {
        var htmlCatProperties = buildHTMLCatProperties(response);

        if (htmlCatProperties) {
            $wrapper.append($(htmlCatProperties));
            $select.after($wrapper);
        }

        if (cb) {
            cb();
        }

    }).fail(function (response) {
        alert("Status: " + response.status + "; " + response.responseText);

    }).always(function () {
        $select.attr("disabled", false);
    });
}

function buildHTMLCatProperties(oCatData) {
    var reciever = [];

    for (var i = 0; i < oCatData.properties.length; i++) {
        var property = oCatData.properties[i];
        var symbolRequire = property.propertyIsRequire ? ' *' : '';
        var title = property.title;
        var tag = getHTMLTagCatProperty(property);
        var row = [
            '<div class="form__row">',
            '   <div class="form__title">' + title + symbolRequire + '</div>',
            tag,
            '</div>',
        ].join('');

        reciever.push(row);
    }

    return reciever.join('');
}

function getHTMLTagCatProperty(property) {
    var propRequire = property.propertyIsRequire ? 'required="required"' : "";
    var kind = property.kindPropertyName;
    var name = property.name;
    var pos = property.propertyPos;
    var propId = property.propertyId;
    var result = 'unknown';
    var el = '';

    if (kind === 'input') {
        el += '<input name="' + name + '" type="text" ' + propRequire + ' data-pos="' + pos + '" value=""/>';

    } else if (kind === 'input_number') {
        el += '<input name="' + name + '" type="text" ' + propRequire + ' data-pos="' + pos + '" value=""/>';

    } else if (kind === 'input_date') {
        el += '<input name="' + name + '" type="date" ' + propRequire + ' data-pos="' + pos + '" value=""/>';

    } else if (kind === 'input_datetime') {
        el += '<input name="' + name + '" type="datetime-local" ' + propRequire + ' data-pos="' + pos + '" value=""/>';

    } else if (kind === 'textarea') {
        el += '<textarea name="' + name + '" ' + propRequire + ' data-pos="' + pos + '"></textarea>';

    } else if (kind === 'photo') {
        el += '<div class="photo"></div>';

    } else if (kind === 'radio' && property.values) {
        for (var i = 0; i < property.values.length; i++) {
            var oVal = property.values[i];
            el += [
                '<label>',
                '   <input type="radio" value="' + oVal.valueId + '" name="' + name + '" ' + propRequire + ' data-pos="' + oVal.pos + '"/>',
                oVal.title,
                '</label>'
            ].join('');
        }

    } else if (kind === 'checkbox') {
        if (property.values.length > 1) {
            for (var i = 0; i < property.values.length; i++) {
                var oVal = property.values[i];
                el += [
                    '<label>',
                    '   <input type="checkbox" value="' + oVal.valueId + '" name="' + name + '" ' + propRequire + ' data-pos="' + oVal.pos + '"/>',
                    oVal.title,
                    '</label>'
                ].join('');
            }

        } else {
            el += [
                '<label>',
                '   <input type="checkbox" value="' + propId + '" name="' + name + '" ' + propRequire + ' data-pos="' + property.propertyPos + '"/>',
                property.title,
                '</label>'
            ].join('');
        }

    } else if (kind === 'select' && property.values.length) {
        el += '<select name="' + name + '" ' + propRequire + ' data-pos="' + pos + '">';
        el += '<option selected="selected" value="" data-pos="0"></option>';

        for (var i = 0; i < property.values.length; i++) {
            var oVal = property.values[i];
            el += '<option value="' + oVal.valueId + '" data-pos="' + oVal.pos + '">' + oVal.title + '</option>';
        }

        el += '</select>';
    }

    if (el) {
        result = el;
    }

    return result;
}
