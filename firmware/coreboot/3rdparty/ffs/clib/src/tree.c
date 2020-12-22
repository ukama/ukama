/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/tree.c $                                             */
/*                                                                        */
/* OpenPOWER FFS Project                                                  */
/*                                                                        */
/* Contributors Listed Below - COPYRIGHT 2014,2015                        */
/* [+] International Business Machines Corp.                              */
/*                                                                        */
/*                                                                        */
/* Licensed under the Apache License, Version 2.0 (the "License");        */
/* you may not use this file except in compliance with the License.       */
/* You may obtain a copy of the License at                                */
/*                                                                        */
/*     http://www.apache.org/licenses/LICENSE-2.0                         */
/*                                                                        */
/* Unless required by applicable law or agreed to in writing, software    */
/* distributed under the License is distributed on an "AS IS" BASIS,      */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or        */
/* implied. See the License for the specific language governing           */
/* permissions and limitations under the License.                         */
/*                                                                        */
/* IBM_PROLOG_END_TAG                                                     */

/*
 *   File: tree.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr:
 *   Note:
 *   Date: 08/21/10
 */

#include <unistd.h>
#include <stdarg.h>
#include <stdlib.h>
#include <malloc.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <limits.h>

#include "libclib.h"
#include "tree.h"

/* ======================================================================= */

tree_node_t *tree_node_prev(tree_node_t * self)
{
	assert(self != NULL);

	if (tree_node_left(self) != NULL) {
		self = tree_node_left(self);

		while (tree_node_right(self) != NULL)
			self = tree_node_right(self);
	} else {
		tree_node_t *parent = tree_node_parent(self);

		while (parent != NULL && self == tree_node_left(parent))
			self = parent, parent = tree_node_parent(parent);

		self = parent;
	}

	return self;
}

tree_node_t *tree_node_next(tree_node_t * self)
{
	assert(self != NULL);

	if (tree_node_right(self) != NULL) {
		self = tree_node_right(self);

		while (self != NULL && tree_node_left(self) != NULL)
			self = tree_node_left(self);
	} else {
		tree_node_t *parent = tree_node_parent(self);

		while (parent != NULL && self == tree_node_right(parent))
			self = parent, parent = tree_node_parent(parent);

		self = parent;
	}

	return self;
}

/* ======================================================================= */

int tree_init(tree_t * self, compare_f compare)
{
	assert(self != NULL);
	assert(compare != NULL);

	self->root = NULL;
	self->min = self->max = NULL;
	self->compare = compare;
	self->size = 0;

	return 0;
}

static inline void __tree_new_min(tree_t * self, tree_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	if (self->min == NULL)
		self->min = node;
	else if (self->compare(node->key, self->min->key) < 0)
		self->min = node;
}

static inline void __tree_new_max(tree_t * self, tree_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	if (self->max == NULL)
		self->max = node;
	else if (self->compare(node->key, self->max->key) > 0)
		self->max = node;
}

static inline void
__tree_new_root(tree_t * self, tree_node_t * node, tree_node_t * left,
		tree_node_t * right)
{
	assert(self != NULL);
	assert(node != NULL);

	node->left = left;
	node->right = right;
	node->parent = NULL;

	if (right != NULL)
		right->parent = node;
	if (left != NULL)
		left->parent = node;

	self->root = node;
}

int tree_insert(tree_t * self, tree_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	__tree_new_min(self, node);
	__tree_new_max(self, node);

	if (self->root == NULL) {
		__tree_new_root(self, node, NULL, NULL);
		self->size++;
	} else {
		tree_node_t *root = self->root;

		while (root != NULL) {
			int rc = self->compare(node->key, root->key);

			if (rc < 0) {
				if (root->left == NULL) {
					node->parent = root;
					root->left = node;
					self->size++;
					return 0;
				} else {
					root = root->left;
				}
			} else if (0 < rc) {
				if (root->right == NULL) {
					node->parent = root;
					root->right = node;
					self->size++;
					return 0;
				} else {
					root = root->right;
				}
			} else {
				UNEXPECTED("duplicate key detected during "
					   "insert");
				return -1;
			}
		}
	}

	return 0;
}

static inline tree_node_t *__tree_find_min(tree_node_t * node)
{
	tree_node_t *min = node;

	if (node != NULL && node->right != NULL) {
		min = node->right;

		while (min != NULL && min->left != NULL)
			min = min->left;
	}

	return min;
}

static inline tree_node_t *__tree_find_max(tree_node_t * node)
{
	tree_node_t *max = node;

	if (node != NULL && node->left != NULL) {
		max = node->left;

		while (max != NULL && max->right != NULL)
			max = max->right;
	}

	return max;
}

int tree_remove(tree_t * self, tree_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	if (self->root == NULL || node == NULL)
		return 0;

	/* =========================== */

	inline tree_node_t *__tree_find_min(tree_node_t * node) {
		tree_node_t *min = node;

		if (node != NULL && node->right != NULL) {
			min = node->right;

			while (min != NULL && min->left != NULL)
				min = min->left;
		}

		return min;
	}

	inline tree_node_t *__tree_find_max(tree_node_t * node) {
		tree_node_t *max = node;

		if (node != NULL && node->left != NULL) {
			max = node->left;

			while (max != NULL && max->right != NULL)
				max = max->right;
		}

		return max;
	}

	inline void remove_single_child(tree_node_t * node, bool left) {
		tree_node_t *new_root = node->left;

		if (left == false)
			new_root = node->right;

		if (node->parent != NULL) {
			new_root->parent = node->parent;

			if (node->parent->left == node)	// handle zig-zag
				node->parent->left = new_root;
			else if (node->parent->right == node)
				node->parent->right = new_root;
		}

		if (node == self->root) {
			self->root = new_root;
			self->root->parent = NULL;
		}
	}

	/* =========================== */

	if (node == self->min)
		self->min = tree_node_next(node);
	if (node == self->max)
		self->max = tree_node_prev(node);

	if (tree_node_leaf(node)) {
		if (node->parent != NULL) {
			if (node->parent->left == node)	// left or right child?
				node->parent->left = NULL;
			else if (node->parent->right == node)
				node->parent->right = NULL;
		}
	} else if (tree_node_internal(node)) {	// two children!
		tree_node_t *root = __tree_find_min(node);
		assert(root != NULL);	// must have a right child

		if (node->right == root) {	// new 'root' is largest child
			root->left = node->left;
			if (node->left != NULL)
				root->left->parent = root;
		} else {	// new 'root' is smallest grandchild
			root->parent->left = root->right;
			if (root->right != NULL)
				root->right->parent = root->parent;

			root->right = node->right;
			root->left = node->left;

			node->right->parent = root;
			node->left->parent = root;
		}

		root->parent = node->parent;

		if (self->root == node) {	// find new parent
			self->root = root;
		} else {
			if (node->parent->left == node)
				node->parent->left = root;
			else if (node->parent->right == node)
				node->parent->right = root;
		}
	} else if (node->left != NULL) {	// single left child
		remove_single_child(node, true);
	} else if (node->right != NULL) {	// single right child
		remove_single_child(node, false);
	}

	node->right = node->left = node->parent = NULL;

	if (0 < self->size)
		self->size--;

	if (self->size <= 0)
		self->root = self->min = self->max = NULL;

	return 0;
}

tree_node_t *tree_find(tree_t * self, const void *key)
{
	assert(self != NULL);
	assert(key != NULL);

	tree_node_t *root = self->root;

	while (root != NULL) {
		int rc = self->compare(key, root->key);

		if (rc < 0)
			root = root->left;
		else if (0 < rc)
			root = root->right;
		else
			break;
	}

	return root;
}

int tree_walk(tree_t * self, tree_walk_f walk_func)
{
	assert(self != NULL);
	assert(walk_func != NULL);

	int __tree_walk(tree_node_t * root) {
		int rc = 0;

		if (root != NULL) {
			__tree_walk(root->left);
			rc = walk_func(root);
			__tree_walk(root->right);
		}

		return rc;
	}

	return __tree_walk(self->root);
}

void tree_dump(tree_t * self, FILE * out)
{
	assert(self != NULL);

	int __tree_node_dump(tree_node_t * node) {
		if (node != NULL)
			return fprintf(out, "node[%p] left[%p] right[%p] "
				       "parent[%p] -- key[%ld]\n",
				       node, node->left, node->right,
				       node->parent, (intptr_t) node->key);
		else
			return 0;
	}

	if (out == NULL)
		out = stdout;

	fprintf(out, "root[%p] min[%p] max[%p] compare[%p] size[%d]\n",
		self->root, self->min, self->max, self->compare, self->size);

	tree_walk(self, __tree_node_dump);
}

void tree_node_dump(tree_node_t * node, FILE * out)
{
	if (node == NULL)
		return;

	void __tree_node_dump(tree_node_t * root, int level) {
		if (root != NULL) {
			if (0 < level) {
				for (int i = 0; i < level; i++)
					fprintf(out, "  ");
			}

			fprintf(out, "node:[%p] left[%p] right[%p] parent[%p] "
				"key[%ld]\n", root, root->left, root->right,
				root->parent, (intptr_t) node->key);

			level++;
			__tree_node_dump(root->left, level);
			__tree_node_dump(root->right, level);
			level--;
		}
	}

	if (out == NULL)
		out = stdout;

	__tree_node_dump(node, 0);
}

/* ======================================================================= */

static tree_node_t *splay(tree_node_t * node, const void *key,
			  compare_f compare)
{
	if (node == NULL)
		return node;

	tree_node_t N = { NULL, NULL, NULL, NULL }
	, *l = &N, *r = &N, *y;

	for (;;) {
		int rc = compare(key, node->key);
		if (rc < 0) {
			if (node->left == NULL)
				break;

			/* rotate right */
			rc = compare(key, node->left->key);
			if (rc < 0) {
				y = node->left;
				node->left = y->right;
				if (y->right != NULL)
					y->right->parent = node;
				y->right = node;
				node->parent = y;
				y->parent = node->parent;
				node = y;

				if (node->left == NULL)
					break;
			}

			/* link right */
			r->left = node;
			node->parent = r;

			r = node;
			node = node->left;
		} else if (0 < rc) {
			if (node->right == NULL)
				break;

			/* rotate left */
			rc = compare(key, node->right->key);
			if (0 < rc) {
				y = node->right;
				node->right = y->left;
				if (y->left != NULL)
					y->left->parent = node;
				y->left = node;
				node->parent = y;
				y->parent = node->parent;
				node = y;

				if (node->right == NULL)
					break;
			}

			/* link left */
			l->right = node;
			node->parent = l;

			l = node;
			node = node->right;
		} else {
			break;
		}
	}

	/* assemble */
	l->right = node->left;
	if (node->left != NULL)
		node->left->parent = l;
	r->left = node->right;
	if (node->right != NULL)
		node->right->parent = r;

	node->left = N.right;
	if (N.right != NULL)
		N.right->parent = node;

	node->right = N.left;
	if (N.left != NULL)
		N.left->parent = node;

	node->parent = NULL;

	return node;
}

int splay_insert(tree_t * self, tree_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	__tree_new_min(self, node);
	__tree_new_max(self, node);

	if (self->root == NULL) {
		node->left = node->right = node->parent = NULL;
		self->root = node;
		self->size = 1;
		return 0;
	}

	self->root = splay(self->root, node->key, self->compare);

	int rc = self->compare(node->key, self->root->key);

	if (rc < 0) {
		node->left = self->root->left;
		if (node->left != NULL)
			node->left->parent = node;
		node->right = self->root;
		self->root->left = NULL;
		self->root->parent = node;
	} else if (0 < rc) {
		node->right = self->root->right;
		if (node->right != NULL)
			node->right->parent = node;
		node->left = self->root;
		self->root->right = NULL;
		self->root->parent = node;
	} else {
		UNEXPECTED("duplicate key detected during insert");
		return -1;
	}

	self->size++;
	self->root = node;

	return 0;
}

int splay_remove(tree_t * self, tree_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	if (self->root == NULL || node == NULL)
		return 0;

	if (node == self->min)
		self->min = tree_node_next(node);
	if (node == self->max)
		self->max = tree_node_prev(node);

	self->root = splay(self->root, node->key, self->compare);

	if (node->key == self->root->key) {	/* found it */
		tree_node_t *x;
#if SAVE
		if (self->root->left == NULL) {
			x = self->root->right;
		} else {
			x = splay(self->root->left, node->key, self->compare);
			x->right = self->root->right;
		}
#else
		if (self->root->left != NULL && self->root->right != NULL) {
			if (__builtin_parity(int64_hash1((int64_t) self->root))) {
				x = splay(self->root->left, node->key,
					  self->compare);
				x->right = self->root->right;
			} else {
				x = splay(self->root->right, node->key,
					  self->compare);
				x->left = self->root->left;
			}
		} else if (self->root->left == NULL) {
			x = self->root->right;
		} else {
			x = splay(self->root->left, node->key, self->compare);
			x->right = self->root->right;
		}
#endif

		self->root->left = self->root->right = NULL;
		self->root->parent = NULL;
		self->root = x;

		if (0 < self->size)
			self->size--;
		if (0 < self->size)
			self->root->parent = NULL;
	}

	return 0;
}

tree_node_t *splay_find(tree_t * self, const void *key)
{
	assert(self != NULL);
	assert(key != NULL);

	self->root = splay(self->root, key, self->compare);

	if (self->root != NULL && self->compare(key, self->root->key))
		return NULL;

	return self->root;
}

/* ======================================================================= */
