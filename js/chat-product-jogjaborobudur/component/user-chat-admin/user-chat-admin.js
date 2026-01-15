    let currentMessages = [];
    let currentChatSession = null;
    let activeChatChannel = null;

    const pusher = new Pusher("4281def63ee8450bd642", {
        cluster: "ap1",
        forceTLS: true
    });

    function subscribeChatChannel(token) {
        if (activeChatChannel) {
            pusher.unsubscribe(activeChatChannel.name);
            activeChatChannel = null;
        }

        activeChatChannel = pusher.subscribe(`chat-${token}`);

        activeChatChannel.bind("new-message", function (data) {
            if (!data || !data.Message) return;

            appendIncomingMessage(data);
        });
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
        $('#sendChatBtnAdmin').prop('disabled', isLoading);

        if (isLoading) {
            $('#sendChatBtnAdmin i')
                .removeClass('fa-paper-plane')
                .addClass('fa-spinner fa-spin');
        } else {
            $('#sendChatBtnAdmin i')
                .removeClass('fa-spinner fa-spin')
                .addClass('fa-paper-plane');
        }
    }
    

    function loadChatRoom(token, session, productId) {
        showChatLoading();
        currentMessages = [];

        $.get("https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/user-session", {
            session: session,
            product_id:productId
        }).then(res => {
            currentChatSession = res.session;
            setChatProductHeader(currentChatSession);
        });

        if (!token) return;

        $.get("https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/messages", {
            token,
            limit: 0,
            offset: 0
        })
        .then(res => {
            currentMessages = res.message || [];
            renderChatGroupedByDate(currentMessages);
            subscribeChatChannel(token);
            hideChatLoading();
        })
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

        if (!currentChatSession?.Token) {
            Swal.fire("Error", "Chat session not ready", "error");
            return;
        }

        setSendLoading(true);
        input.val("");

        const payload = {
            token: currentChatSession.Token,
            sender_type: "admin",
            message
        };

        $.ajax({
            url: "https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/send",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify(payload)
        })
        .then(res => {
            // appendIncomingMessage({
            //     Message: message,
            //     SenderType: "admin",
            //     CreatedAt: new Date().toISOString()
            // });
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
    document.activeElement?.blur();

    if (activeChatChannel) {
        pusher.unsubscribe(activeChatChannel.name);
        activeChatChannel = null;
    }

    $("#chatProductWidgetAdmin").fadeOut(150);
    $("#chatRoomWidgetAdmin").fadeIn(150);
});


$(document).ready(function () {
    $('#backBtn').on('click', function() {
        $("#chatProductWidgetAdmin").fadeOut(150);
        $('#chatRoomWidgetAdmin').fadeIn(150);
        window.renderChatList();
    })

    $('#sendChatBtnAdmin').on('click', function () {
        sendChatMessage();
    });

    $('#chatInput').on('keypress', function (e) {
        if (e.which === 13 && !e.shiftKey && !$('#sendChatBtnAdmin').prop('disabled')) {
            e.preventDefault();
            sendChatMessage();
        }
    });

});
