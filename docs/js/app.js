$(document).ready(function(){

  $(document).scroll(function(){
    if($(window).scrollTop() <= '0'){
      $('.navbar-default').addClass('nav-top');
    } else {
      $('.navbar-default').removeClass('nav-top');
    }

  });

});