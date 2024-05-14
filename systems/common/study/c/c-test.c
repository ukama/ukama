#include <stdio.h>
#include <stdlib.h>
#include <string.h>



typedef struct __node
{
    int data;
    struct __node *link;
} Node;

Node *HEAD = NULL;

Node *new_node(int data)
{

    Node *node = malloc(sizeof(Node));
    if (node)
    {
        memset(node, '\0', sizeof(Node));
        node->data = data;
        node->link = NULL;
    }
    else
    {
        printf("Error:: Memory allocation failure for Node.\n");
    }

    return node;
}

/* Add new node */
Node *append_node(Node *root, int data)
{

    Node *node = new_node(data);
    if (!node)
    {
        return root;
    }

    /* First Node of the List */
    if (root == NULL)
    {
        root = node;
    }
    else
    {
        /* Iterate to last node */
        Node *iter = root;
        while (iter->link != NULL)
        {
            iter = iter->link;
        }

        iter->link = node;
    }

    return root;
}

/* Remove node */
Node *remove_node(Node *root, int data)
{
    if (!root)
    {
        printf("List is empty.\n");
    }

    Node *prev = NULL;
    Node *iter = root;

    while (iter != NULL)
    {
        if (iter->data != data)
        {
            prev = iter;
            iter = iter->link;
        }
        else
        {
            /* First Node of the List */
            if (iter == root)
            {
                root = root->link;
            }
            else
            {
                prev->link = iter->link;
            }
            printf("Removing element %d.\n", iter->data);
            /* Free */
            free(iter);
            iter = NULL;
        }
    }

    return root;
}

void print_list(Node *root)
{

    Node *head = root;
    printf("\n*****************************************\n");
    while (head != NULL)
    {
        if (head == root)
        {
            printf(" |HEAD| -> %d ", head->data);
        }
        else if (head->link != NULL)
        {
            printf("-> %d ", head->data);
        }
        else
        {
            printf("-> %d -> END", head->data);
        }
        head = head->link;
    }
    printf("\n*****************************************\n");
}

Node *reverse_list(Node *root)
{

    if (!root)
    {
        return root;
    }

    Node *prev = NULL;
    Node *next = NULL;
    Node *curr = root;

    while (curr != NULL)
    {
        next = curr->link;
        curr->link = prev;
        prev = curr;
        curr = next;
    }

    return prev;
}

Node *rec_rev_list(Node *root)
{
	Node* head = NULL;
    if (root->link == NULL)
    {
       head = root;
    }
    else
    {
        head = rec_rev_list(root->link);
        root->link->link = root;
        root->link = NULL;
    }
    return head;
}



int main()
{
    Node *head = NULL;

    head = append_node(head, 1);
    print_list(head);
    head = remove_node(head, 1);
    print_list(head);
    head = append_node(head, 1);
    print_list(head);
    head = append_node(head, 2);
    print_list(head);
    head = append_node(head, 3);
    print_list(head);
    head = append_node(head, 4);
    print_list(head);
    head = append_node(head, 5);
    print_list(head);
    head = append_node(head, 6);
    print_list(head);
    head = remove_node(head, 6);
    print_list(head);
    head = append_node(head, 6);
    print_list(head);

    head = remove_node(head, 3);
    print_list(head);

    head = reverse_list(head);
    print_list(head);

    head = rec_rev_list(head);
    print_list(head);
}
