/**
 * Custom JS for tbot
 * Copyright 2016 - Nathan Osman
 */

$(function() {

    // If the page URL matches one of the navbar links, set it as active
    $('.nav li a[href="' + location.pathname + '"]').parent().addClass('active');
});
