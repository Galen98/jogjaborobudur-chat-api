$(document).ready(function() {
    function generateRandomString(length) {
        const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        let result = "";
        for (let i = 0; i < length; i++) {
            result += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        return result;
    }

    function setLoading(isLoading) {
        if (isLoading) {
            $('#startChatBtn').attr('disabled', true);
            $('.btn-text').addClass('d-none');
            $('.btn-loading').removeClass('d-none');
            $('#fullName, #email').attr('readonly', true);
        } else {
            $('#startChatBtn').attr('disabled', false);
            $('.btn-text').removeClass('d-none');
            $('.btn-loading').addClass('d-none');
            $('#fullName, #email').attr('readonly', false);
        }
    }

    $('#startChatBtn').on('click', function () {

        let fullName = $('#fullName').val();
        let email = $('#email').val();
        let productId = $('#productId').val();
        let thumbnail = $('#thumbnailUrl').val();
        let productName = $('#productName').val();

        if (!fullName || !email) {
            Swal.fire("Error", "Full name & email required", "error");
            return;
        }

        const userSession = generateRandomString(50);
        const chatToken = generateRandomString(50);

        setLoading(true);

        $.ajax({
            url: "https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/user-chat",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                full_name: fullName,
                email: email,
                session: userSession
            })
        })
        .then(() => {
            localStorage.setItem("userSession", userSession);

            return $.ajax({
                url: "https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/init-session",
                method: "POST",
                contentType: "application/json",
                data: JSON.stringify({
                    token: chatToken,
                    session: userSession,
                    product_id: Number(productId),
                    thumbnail: thumbnail,
                    product_name: productName
                })
            });
        })

        .then(() => {
            return $.ajax({
                url: "https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/send",
                method: "POST",
                contentType: "application/json",
                data: JSON.stringify({
                    token: chatToken,
                    sender_type: "admin",
                    message: "Hi! I'm here to help you with information about this tour. Feel free to ask anything!"
                })
            });
        })

        .then(() => {
            $("#formChatWidget").fadeOut(150);
            $("#chatProductWidget").fadeIn(150);
            $("#chatLauncher").fadeOut(150)
            $.get("https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/user-session", {
                session: userSession,
                product_id: Number(productId)
            })
            .then(res => {
                currentChatSession = res.session;
                setChatProductHeader(currentChatSession);
                window.loadChatRoom(currentChatSession.Token);
            })
            setLoading(false);
        })
});

    $(document).on("click", "#chatCloseBtn", function () {
        $("#formChatWidget").fadeOut(150);
        $("#chatLauncher").fadeIn(150);
        $("#chatOverlayMobile").fadeOut(150);
    });
})