/**
 * Custom JS for tbot
 * Copyright 2016 - Nathan Osman
 */

$(function() {

    // If the page URL matches one of the navbar links, set it as active
    $('.nav li a[href="' + location.pathname + '"]').parent().addClass('active');

    // If a counter is on the page, use it
    $('.counter').each(function() {
        var $counter = $(this),
            $input = $($counter.data('for'));
        function update() {
            $counter.text((140 - $input.val().length) + ' chars remaining');
        }
        $input.keyup(update);
        update();
    });
});
