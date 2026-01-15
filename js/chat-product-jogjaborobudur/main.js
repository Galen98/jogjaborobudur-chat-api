(function () {

    const BASE_PATH = "js/chat-product-jogjaborobudur/component"; 

    loadComponent("form-chat");
    loadComponent("room-chat");
    loadComponent("product-chat");
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
            <button id="chatLauncher" 
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
                <i class="fa-solid fa-comments"></i> CHAT THIS TOUR
            </button>
            </div>
        `);

    $("body").append(`
    <div id="mobileActionBar" class="d-lg-none">
        <button id="chatLauncherMobile" class="btn btn-chat-mbl">
            <i class="fa-solid fa-comments"></i>
            CHAT THIS TOUR
        </button>
        <a id="checkAvailabilityMobile" class="btn btn-check-mbl">
            CHECK AVAILABILITY
        </a>
    </div>
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

function validateUserSession(session) {
    return $.ajax({
        url: "https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/user-chat-expired",
        method: "POST",
        contentType: "application/json",
        data: JSON.stringify({
            session: session
        })
    });
}


$(document).on("click", "#chatLauncher, #chatLauncherMobile", async function () { 
    if (window.innerWidth <= 768) {
        $("#chatOverlayMobile").fadeIn(150);
    }
    let userSession = localStorage.getItem("userSession");
    let productId = $('#productId').val()

    $("#formChatWidget").hide();
    $("#chatProductWidget").hide();

    
    if (!userSession || userSession.trim() === "") {
        $("#formChatWidget").fadeIn(150);
        return; 
    }

    
    try {
        const res = await validateUserSession(userSession);
        if (res.expired == true) {
            localStorage.removeItem("userSession");
            $("#formChatWidget").fadeIn(150);
            return;
        }
        $("#chatProductWidget").fadeIn(150);

        $.get("https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/user-session", {
            session: userSession,
            product_id: productId
        })
        .then(res => {
            currentChatSession = res.session;
            setChatProductHeader(currentChatSession);
            window.loadChatRoom(currentChatSession.Token);
        })

    } catch (err) {
        console.error("validateUserSession error:", err);
        //$("#formChatWidget").fadeIn(150);
    }

    if (window.innerWidth <= 768) {
        $(this).fadeIn(150);
    } else {
        $(this).fadeOut(150);
    }
});

