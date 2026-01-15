(function () {
    const BASE_PATH = "js/chat-product-jogjaborobudur/component"; 

    loadComponent("room-chat-admin");
    loadComponent("user-chat-admin");
    loadLauncherButton();
    
    function loadComponent(name) {
        $.get(`${BASE_PATH}/${name}/${name}.html`, function (html) {
            if (!$(`link[href='${BASE_PATH}/${name}/${name}.css']`).length) {
                $("head").append(
                    `<link rel="stylesheet" href="${BASE_PATH}/${name}/${name}.css">`
                );
            }

            if (!$(`link[href='${BASE_PATH}/${name}/${name}.js']`).length) {
                $("head").append(
                    `<script src="${BASE_PATH}/${name}/${name}.js"></script>`
                );
            }
            $("body").append(html);
        });
    }

    function loadLauncherButton() {
        $("body").append(`
            <div id="chatOverlayMobile"></div>
            <div class="d-none d-lg-block">
            <button id="chatLauncherAdmin" 
                class="btn rounded shadow font-jogjaborobudur"
                style="
                    position: fixed;
                    bottom: 20px;
                    right: 20px;
                    z-index: 9999;
                    background-color:#25D366;
                    color:white;
                    font-size:16px;
                ">
                <i class="fa-solid fa-comments"></i> USER CHAT LIST
            </button>
            </div>
        `);

        $("body").append(`
        <button id="chatLauncherAdminMobile" 
            class="btn rounded-circle shadow font-jogjaborobudur d-lg-none"
            style="
                position: fixed;
                bottom: 20px;
                right: 20px;
                width: 60px;
                height: 60px;
                z-index: 9999;
                background-color:#25D366;
                color:white;
                display:flex;
                align-items:center;
                justify-content:center;
                font-size:18px;
            ">
            <i class="fa-solid fa-comments"></i>
        </button>
    `);
    }
})();

function loadAssets() {
    const assets = [
        "https://fonts.cdnfonts.com/css/gt-eesti-display-trial",
        "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css",
        "js/chat-product-jogjaborobudur/chat-widget.css"
    ];

    const jsFiles = [
        "https://cdn.jsdelivr.net/npm/sweetalert2@11"
    ]
    assets.forEach(url => {
        const link = document.createElement("link");
        link.rel = "stylesheet";
        link.href = url;
        document.head.appendChild(link);
    });

    jsFiles.forEach(url => {
        const script = document.createElement("script");
        script.src = url;
        document.head.appendChild(script);
    });
}

loadAssets();


$(document).on("click", "#chatLauncherAdmin, #chatLauncherAdminMobile", async function () { 
    if (window.innerWidth <= 768) {
        $("#chatOverlayMobileAdmin").fadeIn(150);
    }

    //$("#chatProductWidgetAdmin").hide();

    // $('#chatRoomWidgetAdmin').fadeIn(150);
    $('#chatRoomWidgetAdmin').fadeIn(150);
    window.renderChatListAdmin();
    $(this).fadeOut(150);
});

