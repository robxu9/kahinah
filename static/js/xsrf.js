var ajax = $.ajax;
$.extend({
    ajax: function(url, options) {
        if (typeof url === 'object') {
            options = url;
            url = undefined;
        }
        options = options || {};
        url = options.url;
        var xsrftoken = $('meta[name=_xsrf]').attr('content');
        var headers = options.headers || {};
        var domain = document.domain.replace(/\./ig, '\\.');
        if (!/^(http:|https:).*/.test(url) || eval(
                '/^(http:|https:)\\/\\/(.+\\.)*' + domain + '.*/').test(
                url)) {
            headers = $.extend(headers, {
                'X-CSRF-Token': xsrftoken
            });
        }
        options.headers = headers;
        return ajax(url, options);
    }
});
