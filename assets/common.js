async function delay(duration) {
  return new Promise((resolve) => {
    setTimeout(resolve, duration);
  });
}

async function disappear(el) {
  if (!el.parentNode) {
    return;
  }

  (async () => {
    el.classList.add("vanish");
    await delay(1000);

    if (el.parentNode) {
      el.parentNode.removeChild(el);
    }
  })();
}

function showNotification(msg, className, duration) {
  var el = document.createElement("div");
  el.setAttribute("class", className);
  el.setAttribute("role", "alert");
  var cl = document.createElement("span");
  cl.innerHTML = "";
  cl.setAttribute("onclick", "disappear(this.parentNode);");
  cl.classList.add(
    "icon-baseline",
    "icon-after",
    "icon-solid",
    "icon-square-xmark",
    "float-right",
    "text-2xl",
    "ml-1",
    "mb-1",
    "cursor-pointer",
  );

  var sp = document.createElement("span");
  sp.innerHTML = msg;

  el.appendChild(cl);
  el.appendChild(sp);

  setTimeout(function () {
    disappear(el);
  }, duration);

  const notifications = document.getElementById("notifications");

  notifications.appendChild(el);
}

function dropHandler(ev) {
  ev.target.classList.remove("highlight");
  // Prevent default behavior (Prevent file from being opened)
  ev.preventDefault();
  ev.stopPropagation();

  if (!ev.dataTransfer) {
    return;
  }

  const files = ev.dataTransfer.files;
  if (files.length == 0) {
    return;
  }

  const fileInput = document.getElementById("event-upload");
  fileInput.files = files;

  htmx.trigger("#event-upload", "drop");
}

function dragOverHandler(ev) {
  ev.preventDefault();
  ev.stopPropagation();
}
function dragEnterHandler(ev) {
  ev.target.classList.add("highlight");
  ev.preventDefault();
  ev.stopPropagation();
}
function dragLeaveHandler(ev) {
  ev.target.classList.remove("highlight");
  ev.preventDefault();
  ev.stopPropagation();
}
