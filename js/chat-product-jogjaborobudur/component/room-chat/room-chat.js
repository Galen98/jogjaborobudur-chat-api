$(document).ready(function(){
    let chatList = [];
    const userSession = localStorage.getItem("userSession");

    function formatDateTime(datetime) {
        if (!datetime) return "";

        const d = new Date(datetime);

        const day = String(d.getDate()).padStart(2, "0");
        const month = String(d.getMonth() + 1).padStart(2, "0");
        const year = d.getFullYear();

        const hours = String(d.getHours()).padStart(2, "0");
        const minutes = String(d.getMinutes()).padStart(2, "0");

        return `${day}-${month}-${year} ${hours}:${minutes}`;
    }

   function initPusher() {
        if (!window.pusher) {
            Pusher.logToConsole = true;

            window.pusher = new Pusher("4281def63ee8450bd642", {
                cluster: "ap1",
                forceTLS: true,
            });
        }
    }

    let sessionChannel = null;

    function connectChatListPusher(userSession) {
        initPusher();

        if (sessionChannel) {
            sessionChannel.unbind_all();
            window.pusher.unsubscribe("session-" + userSession);
        }

        sessionChannel = window.pusher.subscribe("session-" + userSession);

        sessionChannel.bind("session-update", function (msg) {
            const activeToken = localStorage.getItem("activeChatToken");
            const isActiveChat = activeToken === msg.Token;
            const idx = chatList.findIndex(s => s.Token === msg.Token);

            if (idx !== -1) {
                chatList[idx].UpdatedAt = msg.UpdatedAt;
                chatList[idx].IsRead = msg.IsRead;

                const updated = chatList.splice(idx, 1)[0];
                chatList.unshift(updated);
            } else {
                chatList.unshift({
                    Token: msg.Token,
                    ProductID: msg.ProductID,
                    ProductName: msg.ProductName,
                    Thumbnail: msg.Thumbnail,
                    UpdatedAt: msg.UpdatedAt,
                    IsRead: msg.IsRead
                });
            }

            drawChatList(chatList);
        });
    }

   window.renderChatList = function () {
        const userSession = localStorage.getItem("userSession");
        if (!userSession) return;

        $("#chatList").html(`<p class="text-center">Loading...</p>`);

        $.get("https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/user/chat-session", {
            session: userSession
        })
        .then(res => {
            chatList = res.session;
            drawChatList(chatList);
        })
        .catch(() => {
            $("#chatList").html(`<p class="text-danger">Failed load chat</p>`);
        });
    };

    function drawChatList(list) {
        const container = $("#chatList");
        container.html("");

        if (list.length === 0) {
            container.html(`<p class="text-center text-muted">No chat yet</p>`);
            return;
        }

        list.forEach(item => {
            const badge = !item.IsRead
                ? `<span class="badge-new">New</span>`
                : "";

            container.append(`
                <div class="chat-item ${item.IsRead ? '' : 'unread'}"
                     data-token="${item.Token}"
                     data-product="${item.ProductID}">
                    <img src="${item.Thumbnail}">
                    <div class="chat-item-title">
                        <div class="chat-title-row">
                            ${item.ProductName}
                            ${badge}
                        </div>
                        <p class="text-secondary" style="font-size:12px;">
                          Last update: ${formatDateTime(item.UpdatedAt)}
                        </p>
                    </div>
                </div>
                <hr>
            `);
        });
    }

    if (userSession) {
        connectChatListPusher(userSession);
        window.renderChatList(); 
    }

    $(document).on("click", ".chat-item", function () {
        const token = $(this).data("token");
        const $item = $(this);
        $item.removeClass("unread");
        $item.find(".badge-new").remove();

        $.ajax({
            url: `https://jogjaborobudur-chat-api-galen981806-zgwo876o.leapcell.dev/api/chat/open/user/${token}/user`,
            method: "PATCH"
        })
        .always(() => {
            localStorage.setItem("activeChatToken", token);
            $("#chatRoomWidget").fadeOut(150);
            $("#chatProductWidget").fadeIn(150);

            window.loadChatRoom(token);
        });
    });

    $(document).on("click", "#chatCloseBtnRoom", function () {
        localStorage.removeItem("activeChatToken");
        $("#chatRoomWidget").fadeOut(150);
        $("#chatLauncher").fadeIn(150);
        $("#chatLauncherMobile").fadeIn(150);
        $("#chatOverlayMobile").fadeOut(150);
    });
})