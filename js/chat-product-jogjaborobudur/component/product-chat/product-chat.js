    let currentChatSession = null;
    let currentMessages = [];
    let isLoadingChat = false;

    function initPusher() {
        if (!window.pusher) {
            Pusher.logToConsole = true;
            window.pusher = new Pusher("4281def63ee8450bd642", {
                cluster: "ap1",
                forceTLS: true,
            });
        }
    }

    function showChatLoading() {
        $("#chatLoading").fadeIn(150);
    }

    function hideChatLoading() {
        $("#chatLoading").fadeOut(150);
    }

    function appendIncomingMessage(msg) {
        currentMessages.push(msg);
        renderChatGroupedByDate(currentMessages);
    }

    function escapeHTML(text) {
        return $("<div>").text(text).html();
    }

    let chatChannel = null;

    function connectPusherChat(token) {
        initPusher();

        const channelName = "chat-" + token;

        if (chatChannel && chatChannel.name === channelName) {
            return;
        }

        if (chatChannel) {
            chatChannel.unbind_all();
            window.pusher.unsubscribe(chatChannel.name);
            chatChannel = null;
        }

        chatChannel = window.pusher.subscribe(channelName);

        chatChannel.bind("new-message", function (data) {
            appendIncomingMessage(data);
        });
    }


    function connectSessionPusher(userSession) {
        initPusher();

        const sessionChannel = window.pusher.subscribe("session-" + userSession);

        sessionChannel.bind("updated", function (data) {
            window.renderChatList();
        });
    }


    function formatTime(datetime) {
        const d = new Date(datetime.replace(" ", "T"));
        return d.toLocaleTimeString("id-ID", {
            hour: "2-digit",
            minute: "2-digit"
        });
    }

    function formatDate(datetime) {
        const d = new Date(datetime.replace(" ", "T"));
        return d.toLocaleDateString("id-ID", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric"
        });
    }

    function setChatProductHeader(session) {
        if (!session) return;

        $(".chat-product-img").attr("src", session.Thumbnail);
        $(".product-name").text(session.ProductName);
    }

    function setSendLoading(isLoading) {
        $('#chatInput').prop('disabled', isLoading);
        $('#sendChatBtn').prop('disabled', isLoading);

        if (isLoading) {
            $('#sendChatBtn i')
                .removeClass('fa-paper-plane')
                .addClass('fa-spinner fa-spin');
        } else {
            $('#sendChatBtn i')
                .removeClass('fa-spinner fa-spin')
                .addClass('fa-paper-plane');
        }
    }


    // $("#formChatWidget").fadeOut(150);
    // $("#chatProductWidget").fadeIn(150);

    function loadChatRoom(token) { 
        showChatLoading();
        const userSession = localStorage.getItem("userSession");
        const productId = $('#productId').val();
        if (token) {
            currentChatSession = { Token: token }; 

            $.get("https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/messages", {
                token: token,
                limit: 0,
                offset: 0
            })
            .then(res => {
                currentMessages = res.message;
                renderChatGroupedByDate(res.message);
                hideChatLoading();
                connectPusherChat(token);
            })
            .catch(() => hideChatLoading());

            return;
        }
    }


    function renderChatGroupedByDate(messages) {
        const container = $("#chatMessages");
        container.html("");

        let lastDate = null;

        messages.forEach(msg => {
            const datetime = msg.CreatedAt || msg.Time;
            const dateLabel = formatDate(datetime);
            const timeLabel = formatTime(datetime);

            if (dateLabel !== lastDate) {
                container.append(`
                    <div class="chat-date-separator">
                        ${dateLabel}
                    </div>
                `);
                lastDate = dateLabel;
            }

            const bubbleClass = msg.SenderType === "admin" ? "admin" : "user";

            container.append(`
            <div class="chat-bubble ${bubbleClass}">
                    <div class="chat-text">${escapeHTML(msg.Message)}</div>
                <span class="chat-time">${timeLabel}</span>
            </div>
            `);
        });

        container.scrollTop(container.prop("scrollHeight"));
    }

    function sendChatMessage() {
        const input = $('#chatInput');
        const message = input.val().trim();

        if (!message) return;
        if (!currentChatSession || !currentChatSession.Token) {
            Swal.fire("Error", "Chat session not ready", "error");
            return;
        }

        setSendLoading(true);

        const payload = {
            token: currentChatSession.Token,
            sender_type: "user",
            message: message
        };

        input.val("");

        $.ajax({
            url: "https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/send",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify(payload)
        })
        .fail(() => {
            Swal.fire("Error", "Failed to send message", "error");
        })
        .always(() => {
            setSendLoading(false);
            input.focus();
        });
    }



    window.loadChatRoom = loadChatRoom;
    
$(document).on("click", "#chatCloseBtnProd", function () {

    if (chatChannel) {
        chatChannel.unbind_all();
        window.pusher.unsubscribe(chatChannel.name);
        chatChannel = null;
    }

    currentChatSession = null;
    currentMessages = [];

    $("#chatMessages").html("");
    $("#chatProductWidget").fadeOut(150);
    $("#formChatWidget").fadeOut(150);
    $("#chatLauncher").fadeIn(150);
    $("#chatLauncherMobile").fadeIn(150);
    $("#chatOverlayMobile").fadeOut(150);
});


$(document).ready(function () {
    const userSession = localStorage.getItem("userSession");
    if (userSession) {
        connectSessionPusher(userSession);
    }

    $('#backBtn').on('click', function() {
        $("#chatProductWidget").fadeOut(150);
        $('#chatRoomWidget').fadeIn(150);
        window.renderChatList();
    })

    $('#sendChatBtn').on('click', function () {
        sendChatMessage();
    });

    $('#chatInput').on('keypress', function (e) {
        if (e.which === 13 && !e.shiftKey && !$('#sendChatBtn').prop('disabled')) {
            e.preventDefault();
            sendChatMessage();
        }
    });

});
