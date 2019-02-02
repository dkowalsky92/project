var deleteSvg =
  '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" > <path fill="none" d="M0 0h24v24H0V0z"/> <path d="M14.12 10.47L12 12.59l-2.13-2.12-1.41 1.41L10.59 14l-2.12 2.12 1.41 1.41L12 15.41l2.12 2.12 1.41-1.41L13.41 14l2.12-2.12zM15.5 4l-1-1h-5l-1 1H5v2h14V4zM6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM8 9h8v10H8V9z"/>';
var completeSvg =
  '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" > <path fill="none" d="M0 0h24v24H0V0z"/> <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm4.59-12.42L10 14.17l-2.59-2.58L6 13l4 4 8-8z"/>';

document.getElementById("add").addEventListener("click", addTask);

const getTask = id => {
  const requestOptions = {
    method: "GET",
    headers: {
      Accept: "application/json"
    }
  };

  fetch(`http://localhost:5555/task/${id}`, requestOptions)
    .then(response => {
      console.log(response);
      return response.json();
    })
    .then(json => {
      addToList(json["message"], json["id"], task["completed"]);
    })
    .catch(error => {
      console.error(error);
    });
};

const getAllTasks = () => {
  const requestOptions = {
    method: "GET",
    headers: {
      Accept: "application/json"
    }
  };

  fetch(`http://localhost:5555/task`, requestOptions)
    .then(response => {
      return response.json();
    })
    .then(json => {
      json.forEach(task => {
        console.log(json);
        addToList(task["message"], task["id"], task["completed"]);
      });
    })
    .catch(error => {
      console.error(error);
    });
};

function addTask(e) {
  var value = document.getElementById("new-task").value;
  if (value) {
    const requestOptions = {
      method: "POST",
      headers: {
        Accept: "application/json"
      },
      body: JSON.stringify({
        message: value,
        completed: false
      })
    };

    fetch(`http://localhost:5555/task`, requestOptions)
      .then(response => {
        return response.json();
      })
      .then(json => {
        console.log(json);
        addToList(json["message"], json["id"], json["completed"]);
      })
      .catch(error => {
        console.error(error);
      });
  }
}

function removeTask() {
  const requestOptions = {
    method: "DELETE",
    headers: {}
  };
  const id = this.parentNode.parentNode.getAttribute("task_id");
  const element = this.parentNode.parentNode;
  fetch(`http://localhost:5555/task/${id}`, requestOptions)
    .then(() => {
      removeTaskHTML(element);
    })
    .catch(error => {
      console.error(error);
    });
}

function markCompleted(e) {
  const id = this.parentNode.parentNode.getAttribute("task_id");
  const completed = this.parentNode.parentNode.getAttribute("completed");
  const requestOptions = {
    method: "PUT",
    headers: {},
    body: JSON.stringify({
      id: parseInt(id),
      completed: completed == true ? false : true
    })
  };

  const element = this.parentNode.parentNode;
  fetch(`http://localhost:5555/task`, requestOptions)
    .then(() => {
      moveTaskHTML(element);
    })
    .catch(error => {
      console.error(error);
    });
}

const addToList = (name, id, completed) => {
  const task = addTaskHTML(name, id, completed);
  let list = undefined;
  if (!completed) {
    list = document.getElementById("tasks");
  } else {
    list = document.getElementById("finished-tasks");
  }

  if (list.childNodes == null) {
    list.appendChild(task);
  } else {
    list.insertBefore(task, list.childNodes[0]);
  }
};

const addTaskHTML = (name, id, completed) => {
  var task = document.createElement("li");
  task.innerText = name;
  task.setAttribute("task_id", id);
  task.setAttribute("completed", completed);

  var buttons = document.createElement("div");
  buttons.classList.add("complete-delete-container");

  var remove = document.createElement("button");
  remove.classList.add("delete");
  remove.innerHTML = deleteSvg;
  remove.addEventListener("click", removeTask);

  var sep = document.createElement("span");
  sep.classList.add("complete-delete-sep");

  var complete = document.createElement("button");
  complete.classList.add("complete");
  complete.innerHTML = completeSvg;
  complete.addEventListener("click", markCompleted);

  buttons.appendChild(remove);
  buttons.appendChild(sep);
  buttons.appendChild(complete);

  task.appendChild(buttons);

  return task;
};

const removeTaskHTML = target => {
  var parent = target.parentNode;
  parent.removeChild(target);
};

const moveTaskHTML = target => {
  var parent = target.parentNode;
  console.log(parent);
  var id = parent.id;
  console.log(id);

  var list =
    id == "tasks"
      ? document.getElementById("finished-tasks")
      : document.getElementById("tasks");

  parent.removeChild(target);
  if (list.childNodes.length == 0) {
    list.appendChild(target);
  } else {
    list.insertBefore(target, list.childNodes[0]);
  }
};

document.onload = getAllTasks();
